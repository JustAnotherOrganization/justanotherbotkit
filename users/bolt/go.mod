module github.com/justanotherorganization/justanotherbotkit/users/bolt

require (
	github.com/gogo/protobuf v1.1.1 // indirect
	github.com/justanotherorganization/justanotherbotkit/internal v0.0.0
	github.com/justanotherorganization/justanotherbotkit/proto v0.0.0
	github.com/justanotherorganization/justanotherbotkit/users v0.0.0
	github.com/pkg/errors v0.8.0
	github.com/satori/go.uuid v1.2.0
	go.etcd.io/bbolt v1.3.0
	golang.org/x/sys v0.0.0-20181022134430-8a28ead16f52 // indirect
)

replace github.com/justanotherorganization/justanotherbotkit/users => ../

replace github.com/justanotherorganization/justanotherbotkit/internal => ../../internal

replace github.com/justanotherorganization/justanotherbotkit/proto => ../../proto
