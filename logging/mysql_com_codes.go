package logging

//read more client-server protocol from http://dev.mysql.com/doc/internals/en/text-protocol.html
const (
	ComQuit byte = iota + 1 // 1
	ComInitDB // 2
	ComQuery // 3...
	ComFieldList
	ComCreateDB
	ComDropDB
	ComRefresh
	ComShutdown
	ComStatistics
	ComProcessInfo
	ComConnect
	ComProcessKill
	ComDebug
	ComPing
	ComTime
	ComDelayedInsert
	ComChangeUser
	ComBinlogDump
	ComTableDump
	ComConnectOut
	ComRegisterSlave
	ComStmtPrepare
	ComStmtExecute
	ComStmtSendLongData
	ComStmtClose
	ComStmtReset
	ComSetOption
	ComStmtFetch
)
