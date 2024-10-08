default: swag cli
		@echo -e '\t[CB-Spider] build ./bin/cb-spider....'
		@go mod download
		@go mod tidy
		@go build -o bin/cb-spider ./api-runtime
dyna plugin plug dynamic: swag cli
		@echo -e '\t[CB-Spider] build ./bin/cb-spider with plugin mode...'
		@go mod download
	        @go build -tags dyna -o bin/cb-spider-dyna ./api-runtime
		@./build_all_driver_lib.sh;
cc: swag
		@echo -e '\t[CB-Spider] build ./bin/cb-spider-arm for arm...'
	        GOOS=linux GOARCH=arm go build -o cb-spider-arm ./api-runtime
clean clear:
		@echo -e '\t[CB-Spider] cleaning...'
	        @rm -rf bin/cb-spider bin/cb-spider-dyna bin/cb-spider-arm
	        @rm -rf dist-tmp

cli-dist dist-cli: cli
		@echo -e '\t[CB-Spider] tar spctl... to dist'
		@mkdir -p /tmp/spider/dist/conf 
		@cp ./interface/spctl ./interface/spctl.conf /tmp/spider/dist 1> /dev/null
		@cp ./conf/log_conf.yaml /tmp/spider/dist/conf 1> /dev/null
		@mkdir -p ./dist
		@tar -zcvf ./dist/spctl-`(date +%Y.%m.%d.%H)`.tar.gz -C /tmp/spider/dist ./ 1> /dev/null
		@rm -rf /tmp/spider
cli:
		@echo -e '\t[CB-Spider] build ./interface/spctl...'
		@go mod download
		@go mod tidy
		@go build -ldflags="-X 'github.com/cloud-barista/cb-spider/interface/cli/spider/cmd.Version=v0.7.7' \
			-X 'github.com/cloud-barista/cb-spider/interface/cli/spider/cmd.CommitSHA=`(git rev-parse --short HEAD)`' \
			-X 'github.com/cloud-barista/cb-spider/interface/cli/spider/cmd.User=`(id -u -n)`' \
			-X 'github.com/cloud-barista/cb-spider/interface/cli/spider/cmd.Time=`(date)`'" \
			-o ./interface/spctl ./interface/cli/spider/spider.go

swag swagger:
	@echo -e '\t[CB-Spider] generating Swagger documentation'
	@~/go/bin/swag i -g api-runtime/rest-runtime/CBSpiderRuntime.go -d ./,./cloud-control-manager -o api > /dev/null
	@sed -i -e 's/github_com_cloud-barista_cb-spider_cloud-control-manager_cloud-driver_interfaces_resources./spider./g' \
	        -e 's/restruntime./spider./g' ./api/docs.go ./api/swagger.json ./api/swagger.yaml