1. Получение метрик (Prometheus Format)
Здесь выводятся данные из runtime.ReadMemStats, отформатированные для Prometheus.
Запрос:

curl http://localhost:8080/metrics

2. Динамическое управление GC Percent
Меняет параметр debug.SetGCPercent на лету. По умолчанию это 100. Меньшее значение заставляет GC работать чаще (тратит CPU, но экономит RAM). Значение < 0 отключает сборщик мусора.
Запрос:

curl http://localhost:8080/set_gc?val=50
Успешный ответ: GC Percent успешно изменен с 100 на 50

3. Профилирование pprof
Стандартный профайлер Go подключен через импорт _ "net/http/pprof".

Web-интерфейс со списком профилей:
Откройте в браузере: http://localhost:8080/debug/pprof/

Анализ распределения памяти через консоль (требует установленного graphviz для визуализации):

go tool pprof http://localhost:8080/debug/pprof/heap