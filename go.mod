module github.com/micro/micro/v2

go 1.13

require (
	github.com/aws/aws-sdk-go v1.25.27
	github.com/boltdb/bolt v1.3.1
	github.com/challenge-league/nakama-go/commands v0.0.0-00010101000000-000000000000
	github.com/challenge-league/nakama-go/context v0.0.0-00010101000000-000000000000
	github.com/chzyer/logex v1.1.10 // indirect
	github.com/chzyer/readline v0.0.0-20180603132655-2972be24d48e
	github.com/chzyer/test v0.0.0-20180213035817-a1ea475d72b1 // indirect
	github.com/cloudflare/cloudflare-go v0.10.9
	github.com/dustin/go-humanize v1.0.0
	github.com/fsnotify/fsnotify v1.4.7
	github.com/go-acme/lego/v3 v3.4.0
	github.com/golang/protobuf v1.4.2
	github.com/google/uuid v1.1.1
	github.com/gorilla/handlers v1.4.2
	github.com/gorilla/mux v1.7.4
	github.com/hako/durafmt v0.0.0-20200710122514-c0fb7b4da026 // indirect
	github.com/heroiclabs/nakama/v2/apigrpc v0.0.0-00010101000000-000000000000 // indirect
	github.com/micro/cli/v2 v2.1.2
	github.com/micro/go-micro/v2 v2.9.0
	github.com/miekg/dns v1.1.27
	github.com/netdata/go-orchestrator v0.0.0-20190905093727-c793edba0e8f
	github.com/olekukonko/tablewriter v0.0.4
	github.com/patrickmn/go-cache v2.1.0+incompatible
	github.com/pkg/errors v0.9.1
	github.com/serenize/snaker v0.0.0-20171204205717-a683aaf2d516
	github.com/slack-go/slack v0.6.5 // indirect
	github.com/spf13/viper v1.7.0
	github.com/stretchr/testify v1.5.1
	github.com/xlab/treeprint v0.0.0-20181112141820-a009c3971eca
	golang.org/x/crypto v0.0.0-20200510223506-06a226fb4e37
	golang.org/x/net v0.0.0-20200520182314-0ba52f642ac2
	golang.org/x/tools v0.0.0-20191216173652-a0e659d51361
	google.golang.org/genproto v0.0.0-20200226201735-46b91f19d98c
	google.golang.org/grpc v1.27.1
	open-match.dev/open-match v1.0.0 // indirect
)

replace (
	github.com/challenge-league/nakama-go/commands => ./nakama-go/commands
	github.com/challenge-league/nakama-go/context => ./nakama-go/context
	github.com/challenge-league/nakama-go/v2 => ./nakama-go
	github.com/coreos/go-systemd => github.com/coreos/go-systemd/v22 v22.0.0
	github.com/heroiclabs/nakama/v2/apigrpc => ./apigrpc

	github.com/micro/go-micro/v2 => ./go-micro
	google.golang.org/grpc => google.golang.org/grpc v1.26.0
)
