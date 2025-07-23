package orm

func arrayToCommaSeparated(arr []string) string {
	var str string
	for i,v := range arr {
		if(len(arr) == i + 1){
			str += v 	
			break
		}
		str += v + ","
	}
	return str
}

