package define

import "os"

// MailPassword 可以把密码配置在环境变量中（你也可以写在文件中）
//var MailPassword = "MailPassword"//都可以随你喜欢
var MailPassword = os.Getenv("MailPassword")
