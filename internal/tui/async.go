package tui

func runAsync(st *UIState, job func() error, onDone func()) {
	if !st.Busy.CompareAndSwap(false, true) {
		return
	}
	go func() {
		err := job()
		st.App.QueueUpdateDraw(func() {
			st.Busy.Store(false)
			if err != nil {
				showError(st, err)
				return
			}
			if onDone != nil {
				onDone()
			}
		})
	}()
}
