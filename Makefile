app_name = pritunl-http-api

define build-linux
	@rm -f $(app_name)*
        CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -installsuffix cgo
endef


define build-mac
	@rm -f $(app_name)*
	CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build -a -installsuffix cgo
endef


define build-window
	@rm -f $(app_name)*
	CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -a -installsuffix cgo
endef


define stop
	@export pid=`ps -ef | grep $(app_name) | grep -v 'grep' | awk '{print $$2}'`; \
	if [[ "$$pid" != "" ]] ;\
	then \
		echo "kill $$pid ..." ;\
		kill $$pid && sleep 5 && kill -9 $$pid ;\
	fi
endef


define start
	@./$(app_name)
endef


define dev
	@rm -f $(app_name)*
	@go build
endef


stop:
	$(call stop)

start:
	$(call start)

restart:
	$(call stop)
	$(call start)

dev:
	$(call stop)
	$(call dev)
	$(call start)

linux:
	$(call build-linux)

mac:
	$(call build-mac)

window:
	$(call build-window)



.PHONY: stop start dev restart linux mac window