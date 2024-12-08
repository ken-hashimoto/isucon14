include env.sh
# 変数定義 ------------------------

# SERVER_ID: env.sh内で定義

# 問題によって変わる変数
USER:=isucon
BIN_NAME:=isuride
BUILD_DIR:=/home/isucon/webapp/go
SERVICE_NAME:=$(BIN_NAME)-go.service

DB_PATH:=/etc/mysql
NGINX_PATH:=/etc/nginx
SYSTEMD_PATH:=/etc/systemd/system

NGINX_LOG:=/var/log/nginx/access.log
DB_SLOW_LOG:=/var/log/mysql/mysql-slow.log # slow_query_log = OFFとなっている場合、このファイルは存在しない.

# netdataを使う手順: 1. make setup-netdata -> 2. ポートフォワーディング
# netdataを使う場合、ポートフォワーディングをすると便利
# なので使う場合は~/.ssh/configに以下を書くこと
# Host isucon-1
#     HostName <サーバー1のIP>
#     User isucon
#     LocalForward 19999 localhost:19999

# Host isucon-2
#     HostName <サーバー2のIP>
#     User isucon
#     LocalForward 19998 localhost:19999

# Host isucon-3
#     HostName <サーバー3のIP>
#     User isucon
#     LocalForward 19997 localhost:19999

# http://localhost:19999/netdata.confのdirectories.webで確認可能
NETDATA_WEBROOT_PATH:=/var/lib/netdata/www
NETDATA_CUSTUM_HTML:=tool-config/netdata/*

DISCOCAT_TRIPLE_BACK_QUOTES_DIR:=tool-config/discocat
DISCOCAT_TRIPLE_BACK_QUOTES:=tool-config/discocat/triple-back-quotes.txt
DISCOCAT_TMPFILE:=tmp/discocat

# メインで使うコマンド ------------------------


# 一番初めに実行する. これからの処理で必要なファイルがあるかをチェックする. ない場合は上の変数を書き換えること. 
.PHONY: check
check: check-paths check-discocat-config 

# サーバーの環境構築　ツールのインストール、gitまわりのセットアップ
.PHONY: setup
setup: setup-discocat-config install-tools git-setup set-gopath

# 設定ファイルなどを取得してgit管理下に配置する
.PHONY: get-conf
get-conf: check-server-id get-db-conf get-nginx-conf get-service-file get-envsh

# リポジトリ内の設定ファイルをそれぞれ配置する
.PHONY: deploy-conf
deploy-conf: check-server-id deploy-db-conf deploy-nginx-conf deploy-service-file deploy-envsh

# ベンチマークを走らせる直前に実行する
.PHONY: bench
bench: check-server-id discocat-now-status pull-main mv-logs build deploy-conf restart watch-service-log

# slow-queryを実行し、その結果をdiscordに送信する.
.PHONY: discocat-slow-query
discocat-slow-query:
	export PATH=$PATH:$(go env GOPATH)/bin
	@make refresh-descocat-tmp
	echo "SERVER_ID: $(SERVER_ID)" >> $(DISCOCAT_TMPFILE)
	echo "" >> $(DISCOCAT_TMPFILE)
	@make slow-query >> $(DISCOCAT_TMPFILE)
	cat $(DISCOCAT_TMPFILE) | discocat

# alpでアクセスログを確認する. ~/tool-config/alp/config.ymlを準備してから実行すること.
.PHONY: alp
alp:
	sudo alp ltsv --file=$(NGINX_LOG) --config=/home/isucon/tool-config/alp/config.yml

setup-discocat-config:
	mkdir -p ${DISCOCAT_TRIPLE_BACK_QUOTES_DIR}
	echo '```' > ${DISCOCAT_TRIPLE_BACK_QUOTES}

.PHONY: discocat-alp
discocat-alp:
	@make refresh-descocat-tmp
	cat $(DISCOCAT_TRIPLE_BACK_QUOTES) >> $(DISCOCAT_TMPFILE)
	echo "" >> $(DISCOCAT_TMPFILE)
	echo "SERVER_ID: $(SERVER_ID)" >> $(DISCOCAT_TMPFILE)
	echo "" >> $(DISCOCAT_TMPFILE)
	@make alp >> $(DISCOCAT_TMPFILE)
	cat $(DISCOCAT_TRIPLE_BACK_QUOTES) >> $(DISCOCAT_TMPFILE)
	cat $(DISCOCAT_TMPFILE) | discocat

.PHONY: check-paths
check-paths:
	@echo "Checking required paths..."
	@test -d $(BUILD_DIR) || (echo "Error: $(BUILD_DIR) does not exist" && exit 1)
	@test -d $(DB_PATH) || (echo "Error: $(DB_PATH) does not exist" && exit 1)
	@test -d $(NGINX_PATH) || (echo "Error: $(NGINX_PATH) does not exist" && exit 1)
	@test -d $(SYSTEMD_PATH) || (echo "Error: $(SYSTEMD_PATH) does not exist" && exit 1)
	@test -f $(NGINX_LOG) || (echo "Warning: $(NGINX_LOG) does not exist")
	@test -f $(DB_SLOW_LOG) || (echo "Warning: $(DB_SLOW_LOG) does not exist")
	@echo "Path check completed"

# pprofを使う手順
# 0. ポートフォワーディングの設定を追加
#   ssh -L 8090:localhost:8090 isucon@<サーバーIP>
# 1. _ "net/http/pprof" をmain.goのimportに追加
# 2. main.goの最初に以下のコードを追加して、デバッグ用のエンドポイントを別のポートで起動
# 		go func() {
#     		http.ListenAndServe(":6060", nil)
# 		}()
# 3. ベンチマーク開始直前にmake pprof-recordを実行
# 4. ベンチマーク終了後、make pprof-checkを実行
# 5. ローカルのPCからlocalhost:8090で結果を確認

# pprofで記録する. 記録時間はsecondsで調節すること.
.PHONY: pprof-record
pprof-record:
	go tool pprof http://localhost:6060/debug/pprof/profile -seconds 60

# pprofで確認する
.PHONY: pprof-check
pprof-check:
	$(eval latest := $(shell ls -rt pprof/ | tail -n 1))
	go tool pprof -http=localhost:8090 pprof/$(latest)

# DBに接続する
.PHONY: access-db
access-db:
	mysql -h $(MYSQL_HOST) -P $(MYSQL_PORT) -u $(MYSQL_USER) -p$(MYSQL_PASS) $(MYSQL_DBNAME)

# 主要コマンドの構成要素 ------------------------

.PHONY: install-tools
install-tools:
	sudo apt update
	sudo apt upgrade
	sudo apt install -y percona-toolkit dstat git unzip snapd graphviz tree

	# alpのインストール
	wget https://github.com/tkuchiki/alp/releases/download/v1.0.9/alp_linux_amd64.zip
	unzip alp_linux_amd64.zip
	sudo install alp /usr/local/bin/alp
	rm alp_linux_amd64.zip alp

	# netdataのインストール
	wget -O /tmp/netdata-kickstart.sh https://my-netdata.io/kickstart.sh && sh /tmp/netdata-kickstart.sh

	# discocatのインストール
	go install github.com/wan-nyan-wan/discocat@latest

.PHONY: set-gopath
set-gopath:
	# .bashrcに追加（修正版）
	echo 'export PATH=$$PATH:$$(go env GOPATH)/bin' >> ~/.bashrc
	# sourceの代わりに現在のシェルにPATHを追加
	export PATH=$$PATH:$$(go env GOPATH)/bin
.PHONY: git-setup
git-setup:
	# git用の設定は適宜変更して良い
	git config --global user.email "isucon@example.com"
	git config --global user.name "isucon"

	# deploykeyの作成
	ssh-keygen -t ed25519

.PHONY: check-discocat-config
check-discocat-config:
	@if [ ! -f .config/discocat.yml ]; then \
		echo "Error: discocat.yml が見つかりません。"; \
		echo "./config/discocat.yml を作成してから再度実行してください。"; \
		exit 1; \
	fi

# slow queryを確認する
.PHONY: slow-query
slow-query:
	sudo pt-query-digest $(DB_SLOW_LOG)

.PHONY: pull-main
pull-main:
	git pull origin main
.PHONY: check-server-id
check-server-id:
ifdef SERVER_ID
	@echo "SERVER_ID=$(SERVER_ID)"
else
	@echo "SERVER_ID is unset"
	@exit 1
endif

.PHONY: set-as-s1
set-as-s1:
	echo "SERVER_ID=s1" >> env.sh

.PHONY: set-as-s2
set-as-s2:
	echo "SERVER_ID=s2" >> env.sh

.PHONY: set-as-s3
set-as-s3:
	echo "SERVER_ID=s3" >> env.sh

.PHONY: get-db-conf
get-db-conf:
	sudo mkdir -p ~/$(SERVER_ID)/etc/mysql
	sudo cp -R $(DB_PATH)/* ~/$(SERVER_ID)/etc/mysql
	sudo chown $(USER) -R ~/$(SERVER_ID)/etc/mysql

.PHONY: get-nginx-conf
get-nginx-conf:
	sudo mkdir -p ~/$(SERVER_ID)/etc/nginx
	sudo cp -R $(NGINX_PATH)/* ~/$(SERVER_ID)/etc/nginx
	sudo chown $(USER) -R ~/$(SERVER_ID)/etc/nginx

.PHONY: get-service-file
get-service-file:
	sudo mkdir -p ~/$(SERVER_ID)/etc/systemd/system
	sudo cp $(SYSTEMD_PATH)/$(SERVICE_NAME) ~/$(SERVER_ID)/etc/systemd/system/$(SERVICE_NAME)
	sudo chown $(USER) ~/$(SERVER_ID)/etc/systemd/system/$(SERVICE_NAME)

.PHONY: get-envsh
get-envsh:
	sudo mkdir -p ~/$(SERVER_ID)/home/isucon
	sudo cp ~/env.sh ~/$(SERVER_ID)/home/isucon/env.sh

.PHONY: deploy-db-conf
deploy-db-conf:
	sudo cp -R ~/$(SERVER_ID)/etc/mysql/* $(DB_PATH)

.PHONY: deploy-nginx-conf
deploy-nginx-conf:
	sudo cp -R ~/$(SERVER_ID)/etc/nginx/* $(NGINX_PATH)

.PHONY: deploy-service-file
deploy-service-file:
	sudo cp ~/$(SERVER_ID)/etc/systemd/system/$(SERVICE_NAME) $(SYSTEMD_PATH)/$(SERVICE_NAME)

.PHONY: deploy-envsh
deploy-envsh:
	sudo cp ~/$(SERVER_ID)/home/isucon/env.sh ~/env.sh

.PHONY: build
build:
	cd $(BUILD_DIR); \
	go build -o $(BIN_NAME)

.PHONY: restart
restart:
	sudo systemctl daemon-reload
	sudo systemctl restart $(SERVICE_NAME)
	sudo systemctl restart mysql
	sudo systemctl restart nginx
	sudo chmod 644 /etc/powerdns/pdns.conf
	sudo systemctl restart pdns.service


.PHONY: mv-logs
mv-logs:
	$(eval when := $(shell date "+%s"))
	sudo mkdir -p ~/logs/$(when)
	sudo test -f $(NGINX_LOG) && \
		sudo mv -f $(NGINX_LOG) ~/logs/nginx/$(when)/ || echo ""
	sudo test -f $(DB_SLOW_LOG) && \
		sudo mv -f $(DB_SLOW_LOG) ~/logs/mysql/$(when)/ || echo ""

.PHONY: watch-service-log
watch-service-log:
	sudo journalctl -u $(SERVICE_NAME) -n10 -f

.PHONY: list-services
list-services:
	@echo "Active services status:"
	@sudo systemctl list-units --type=service --state=running


.PHONY: netdata-setup
netdata-setup:
	sudo cp $(NETDATA_CUSTUM_HTML) $(NETDATA_WEBROOT_PATH)/

.PHONY: $(DISCOCAT_TMPFILE)
refresh-descocat-tmp:
	rm -f $(DISCOCAT_TMPFILE)
	mkdir -p tmp
	touch $(DISCOCAT_TMPFILE)

.PHONY: discocat-now-status
discocat-now-status:
	export PATH=$PATH:$(go env GOPATH)/bin
	@make refresh-descocat-tmp
	echo "----------------------------------------------------------------" >> $(DISCOCAT_TMPFILE)
	cat $(DISCOCAT_TRIPLE_BACK_QUOTES) >> $(DISCOCAT_TMPFILE)
	echo "SERVER_ID: $(SERVER_ID)" >> $(DISCOCAT_TMPFILE)
	git branch --contains=HEAD >> $(DISCOCAT_TMPFILE)
	TZ=JST-9 date >> $(DISCOCAT_TMPFILE)
	echo "" >> $(DISCOCAT_TMPFILE)
	git show -s >> $(DISCOCAT_TMPFILE)
	cat $(DISCOCAT_TRIPLE_BACK_QUOTES) >> $(DISCOCAT_TMPFILE)
	cat $(DISCOCAT_TMPFILE) | discocat

.PHONY: setup-netdata
setup-netdata:
	# netdata.confのバックアップを作成
	sudo cp /etc/netdata/netdata.conf /etc/netdata/netdata.conf.backup
	# 設定ファイルを更新
	echo '[web]' | sudo tee /etc/netdata/netdata.conf
	echo '    bind to = *' | sudo tee -a /etc/netdata/netdata.conf
	# プライベートネットワークからのアクセスのみを許可
	echo '    allow connections from = localhost 127.0.0.1 ::1 192.168.* 10.* 172.16.* 172.17.* 172.18.* 172.19.* 172.20.* 172.21.* 172.22.* 172.23.* 172.24.* 172.25.* 172.26.* 172.27.* 172.28.* 172.29.* 172.30.* 172.31.*' | sudo tee -a /etc/netdata/netdata.conf
	# netdataを再起動
	sudo systemctl restart netdata

.PHONY: restore-netdata
restore-netdata:
	# 設定を元に戻す
	sudo mv /etc/netdata/netdata.conf.backup /etc/netdata/netdata.conf
	sudo systemctl restart netdata


.PHONY: stop-netdata
stop-netdata:
	# netdataサービスを停止
	sudo systemctl stop netdata