package global

const DbString = "root:mysql@tcp(localhost:3306)/napTune?charset=utf8&parseTime=True"

func CheckErr(err error) {
	if err != nil {
		panic(err)
	}
}