package proxy

import (
	"database/sql"
	"fmt"
	"github.com/caybokotze/dbmux/database"
	"github.com/caybokotze/dbmux/logging"
	"log"
	"strconv"
	"strings"
)

type Query struct {
	BindPort   int64
	ClientIP   string
	ClientPort int64
	ServerIP   string
	ServerPort int64
	SqlType    string
	SqlString  string
}

func ipPortFromNetAddr(s string) (ip string, port int64) {
	addrInfo := strings.SplitN(s, ":", 2)
	ip = addrInfo[0]
	port, _ = strconv.ParseInt(addrInfo[1], 10, 64)
	return
}

func convertToUnixLine(sql string) string {
	sql = strings.Replace(sql, "\r\n", "\n", -1)
	sql = strings.Replace(sql, "\r", "\n", -1)
	return sql
}

func sqlEscape(s string) string {
	var j = 0
	if len(s) == 0 {
		return ""
	}

	tempStr := s[:]
	desc := make([]byte, len(tempStr)*2)
	for i := 0; i < len(tempStr); i++ {
		flag := false
		var escape byte
		switch tempStr[i] {
		case '\r':
			flag = true
			escape = '\r'
			break
		case '\n':
			flag = true
			escape = '\n'
			break
		case '\\':
			flag = true
			escape = '\\'
			break
		case '\'':
			flag = true
			escape = '\''
			break
		case '"':
			flag = true
			escape = '"'
			break
		case '\032':
			flag = true
			escape = 'Z'
			break
		default:
		}
		if flag {
			desc[j] = '\\'
			desc[j+1] = escape
			j = j + 2
		} else {
			desc[j] = tempStr[i]
			j = j + 1
		}
	}
	return string(desc[0:j])
}

type LogConfiguration struct {
	Source       *Connection
	Destination  *Connection
	BufferSize   uint
	Verbosity    bool
	DatabaseHost *sql.DB
}

func Log(config LogConfiguration) {
	buffer := make([]byte, config.BufferSize)
	var sqlInfo Query
	sqlInfo.ClientIP, sqlInfo.ClientPort = ipPortFromNetAddr(config.Source.Connection.RemoteAddr().String())
	sqlInfo.ServerIP, sqlInfo.ServerPort = ipPortFromNetAddr(config.Destination.Connection.RemoteAddr().String())
	_, sqlInfo.BindPort = ipPortFromNetAddr(config.Source.Connection.LocalAddr().String())

	for {
		n, err := config.Source.Read(buffer)
		if err != nil {
			return
		}
		if n >= 5 {
			var verboseStr string
			switch buffer[4] {
			case logging.ComQuit:
				verboseStr = fmt.Sprintf("From %s To %s; Quit: %s\n", sqlInfo.ClientIP, sqlInfo.ServerIP, "user quit")
				sqlInfo.SqlType = "Quit"
			case logging.ComInitDB:
				verboseStr = fmt.Sprintf("From %s To %s; schema: use %s\n", sqlInfo.ClientIP, sqlInfo.ServerIP, string(buffer[5:n]))
				sqlInfo.SqlType = "Schema"
			case logging.ComQuery:
				verboseStr = fmt.Sprintf("From %s To %s; Query: %s\n", sqlInfo.ClientIP, sqlInfo.ServerIP, string(buffer[5:n]))
				sqlInfo.SqlType = "Query"
			case logging.ComCreateDB:
				verboseStr = fmt.Sprintf("From %s To %s; CreateDB: %s\n", sqlInfo.ClientIP, sqlInfo.ServerIP, string(buffer[5:n]))
				sqlInfo.SqlType = "CreateDB"
			case logging.ComDropDB:
				verboseStr = fmt.Sprintf("From %s To %s; DropDB: %s\n", sqlInfo.ClientIP, sqlInfo.ServerIP, string(buffer[5:n]))
				sqlInfo.SqlType = "DropDB"
			case logging.ComRefresh:
				verboseStr = fmt.Sprintf("From %s To %s; Refresh: %s\n", sqlInfo.ClientIP, sqlInfo.ServerIP, string(buffer[5:n]))
				sqlInfo.SqlType = "Refresh"
			case logging.ComStmtPrepare:
				verboseStr = fmt.Sprintf("From %s To %s; Prepare Query: %s\n", sqlInfo.ClientIP, sqlInfo.ServerIP, string(buffer[5:n]))
				sqlInfo.SqlType = "Prepare Query"
			case logging.ComStmtExecute:
				verboseStr = fmt.Sprintf("From %s To %s; Prepare Args: %s\n", sqlInfo.ClientIP, sqlInfo.ServerIP, string(buffer[5:n]))
				sqlInfo.SqlType = "Prepare Args"
			case logging.ComProcessKill:
				verboseStr = fmt.Sprintf("From %s To %s; Kill: kill conntion %s\n", sqlInfo.ClientIP, sqlInfo.ServerIP, string(buffer[5:n]))
				sqlInfo.SqlType = "Kill"
			default:
			}

			if config.Verbosity {
				log.Print(verboseStr)
			}

			if strings.EqualFold(sqlInfo.SqlType, "Quit") {
				sqlInfo.SqlString = "user quit"
			} else {
				sqlInfo.SqlString = convertToUnixLine(sqlEscape(string(buffer[5:n])))
			}

			if !strings.EqualFold(sqlInfo.SqlType, "") && config.DatabaseHost != nil {
				insertLog(config.DatabaseHost, &sqlInfo)
			}
		}

		_, err = config.Destination.Write(buffer[0:n])
		if err != nil {
			return
		}
	}
}

func insertLog(db *sql.DB, t *Query) bool {
	insertSql := `
	insert into query_log(bindport, client, client_port, server, server_port, sql_type, 
	sql_string, create_time) values (%d, '%s', %d, '%s', %d, '%s', '%s', now())
	`
	_, err := database.ExecQuery(db, fmt.Sprintf(
		insertSql,
		t.BindPort,
		t.ClientIP,
		t.ClientPort,
		t.ServerIP,
		t.ServerPort,
		t.SqlType,
		t.SqlString))

	if err != nil {
		return false
	}
	return true
}
