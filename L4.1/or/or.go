package or

// Or объединяет произвольное количество done-каналов в один.
// Возвращаемый канал закрывается, как только закрывается любой из входящих каналов.
func Or(channels ...<-chan interface{}) <-chan interface{} {
	// Базовые случаи: если каналов нет или он всего один
	switch len(channels) {
	case 0:
		return nil
	case 1:
		return channels[0]
	}

	orDone := make(chan interface{})

	go func() {
		defer close(orDone)

		switch len(channels) {
		case 2:
			// Если канала два, просто ждем любой из них
			select {
			case <-channels[0]:
			case <-channels[1]:
			}
		default:
			// Оптимизация: проверяем первые 3 канала, а остаток отправляем в рекурсию.
			// КРИТИЧЕСКИ ВАЖНО: мы добавляем orDone в конец среза для рекурсивного вызова.
			// Если сработает channels[0], [1] или [2], orDone закроется через defer,
			// что заставит нижестоящую рекурсивную ветку мгновенно схлопнуться.
			select {
			case <-channels[0]:
			case <-channels[1]:
			case <-channels[2]:
			case <-Or(append(channels[3:], orDone)...):
			}
		}
	}()

	return orDone
}
