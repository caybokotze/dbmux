package logging

import (
	"fmt"
	"github.com/caybokotze/dbmux/database"
	"github.com/caybokotze/dbmux/tcp"
	"log"
	"strconv"
	"strings"
)



type Query struct {
	BindPort  int64
	ClientIP    string
	ClientPort     int64
	ServerIP    string
	ServerPort     int64
	SqlType   string
	SqlString string
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

func ProxyLog(source, destination *tcp.Connection, bufferSize uint) {
	buffer := make([]byte, bufferSize)
	var sqlInfo Query
	sqlInfo.ClientIP, sqlInfo.ClientPort = ipPortFromNetAddr(source.Connection.RemoteAddr().String())
	sqlInfo.ServerIP, sqlInfo.ServerPort = ipPortFromNetAddr(destination.Connection.RemoteAddr().String())
	_, sqlInfo.BindPort = ipPortFromNetAddr(source.Connection.LocalAddr().String())

	for {
		n, err := source.Read(buffer)
		if err != nil {
			return
		}
		if n >= 5 {
			var verboseStr string
			switch buffer[4] {
			case comQuit:
				verboseStr = fmt.Sprintf("From %s To %s; Quit: %s\n", sqlInfo.ClientIP, sqlInfo.ServerIP, "user quit")
				sqlInfo.SqlType = "Quit"
			case comInitDB:
				verboseStr = fmt.Sprintf("From %s To %s; schema: use %s\n", sqlInfo.ClientIP, sqlInfo.ServerIP, string(buffer[5:n]))
				sqlInfo.SqlType = "Schema"
			case comQuery:
				verboseStr = fmt.Sprintf("From %s To %s; Query: %s\n", sqlInfo.ClientIP, sqlInfo.ServerIP, string(buffer[5:n]))
				sqlInfo.SqlType = "Query"
			case comCreateDB:
				verboseStr = fmt.Sprintf("From %s To %s; CreateDB: %s\n", sqlInfo.ClientIP, sqlInfo.ServerIP, string(buffer[5:n]))
				sqlInfo.SqlType = "CreateDB"
			case comDropDB:
				verboseStr = fmt.Sprintf("From %s To %s; DropDB: %s\n", sqlInfo.ClientIP, sqlInfo.ServerIP, string(buffer[5:n]))
				sqlInfo.SqlType = "DropDB"
			case comRefresh:
				verboseStr = fmt.Sprintf("From %s To %s; Refresh: %s\n", sqlInfo.ClientIP, sqlInfo.ServerIP, string(buffer[5:n]))
				sqlInfo.SqlType = "Refresh"
			case comStmtPrepare:
				verboseStr = fmt.Sprintf("From %s To %s; Prepare Query: %s\n", sqlInfo.ClientIP, sqlInfo.ServerIP, string(buffer[5:n]))
				sqlInfo.SqlType = "Prepare Query"
			case comStmtExecute:
				verboseStr = fmt.Sprintf("From %s To %s; Prepare Args: %s\n", sqlInfo.ClientIP, sqlInfo.ServerIP, string(buffer[5:n]))
				sqlInfo.SqlType = "Prepare Args"
			case comProcessKill:
				verboseStr = fmt.Sprintf("From %s To %s; Kill: kill conntion %s\n", sqlInfo.ClientIP, sqlInfo.ServerIP, string(buffer[5:n]))
				sqlInfo.SqlType = "Kill"
			default:
			}

			if sourceverbosityEnabled {
				log.Print(verboseStr)
			}

			if strings.EqualFold(sqlInfo.SqlType, "Quit") {
				sqlInfo.SqlString = "user quit"
			} else {
				sqlInfo.SqlString = convertToUnixLine(sqlEscape(string(buffer[5:n])))
			}

			if !strings.EqualFold(sqlInfo.SqlType, "") && DatabaseHost != nil {
				database.InsertLog(main.DatabaseHost, &sqlInfo)
			}
		}

		_, err = destination.Write(buffer[0:n])
		if err != nil {
			return
		}
	}
}
