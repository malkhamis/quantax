package history

func init() {

	err := validateAllFormulas()
	if err != nil {
		panic(err)
	}
}
