- 数据库名：orbia，表名称 orbia_xxxx ，关于数据库里的处理你可以使用 mysql的 mcp 工具来处理。
- 任何功能的修改优先使用 hz 框架 来生成基本脚手架代码，包括更新 thrift 定义后，或相关代码后，都需要调用 hz 来重新生成脚手架。
- router、handler 代码只包含路由 和基本的处理信息。
- service 目录的代码包含业务逻辑，并且有一个 rpc 目录，比如xxx_rpc.go ，用于抽象一些通用的逻辑，api 的逻辑放在service/xxx_service.go 中。
- 使用 go run . 运行服务。
- 本项目里的所有接口定义都只用 post method，包括 idl 定义，入参出参都用 json。
- 遇到数据库结构的任何问题，应以sql/init.sql 为准，通过执行 script/init_db.sh 来更新数据库。
- dal/mysql 目录下的代码是数据库的模型
- biz/model 目录下的代码是 api 的模型
- 所有接口返回值应该遵循 common 里配置的通用包装，如：
struct BaseResp {
    1: i32 code
    2: string message
}
- 符合 golang 的最佳代码实践以及充分代码鲁棒性和复用性。
- 如果你在回复用户的过程中，发现一个认为可以被设定成通用规则的规则，请写到本文件最后，规则应该简单，一般一句话就能描述。