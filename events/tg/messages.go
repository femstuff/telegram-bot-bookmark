package tg

const helpMsg = `Я могу сохранять твои закладки, чтобы не копить их в браузере.
Также я могу отправлять их тебе для прочтения.

Для того, чтобы сохранить интересующую ссылку - отправь мне ее сообщением.

Для того, чтобы получить случайную ссылку из хранимых - отправь команду /random.
Внимание! После получения ссылки, она будет удалена из списка.`

const helloMsg = "Привет, вот справка по использованию бота: \n" + helpMsg

const (
	unknownMsg       = "Неизвестная команда."
	enoughPagesMsg   = "Нет сохраненных ссылок."
	successfulMsg    = "Ресурс успешно сохранен."
	alreadyExistsMsg = "Данный ресурс уже был сохранен."
)