package orm

func arrayToCommaSeparatedTable(arr []string, table string) string {
	var str string
	for i,v := range arr {
		if(len(arr) == i + 1){
			str += table + "." + v 	
			break
		}
		str += table + "." + v + ","
	}
	return str
}

