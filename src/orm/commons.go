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

func columnArrayToCommaSeparatedTable(arr []column, table string) string {
	var str string
	for i,v := range arr {
		if(len(arr) == i + 1) {
			str += table + "." + v.Name 	
			break
		}
		str += table + "." + v.Name + ","
	}
	return str
}

