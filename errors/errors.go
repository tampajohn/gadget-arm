package errors

func Check(err error, bad ...func(error)) {
	if err != nil {
		if bad != nil {
			for _, f := range bad {
				f(err)
			}
		} else {
			panic(err)
		}
	}
}
