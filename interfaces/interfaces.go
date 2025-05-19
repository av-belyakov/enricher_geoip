package interfaces

//**************** счётчик *****************

type Counter interface {
	SendMessage(string, int)
}

//************** логирование ***************

type Logger interface {
	GetChan() <-chan Messager
	Send(msgType, msgData string)
}

type Messager interface {
	GetType() string
	SetType(v string)
	GetMessage() string
	SetMessage(v string)
}

type WriterLoggingData interface {
	Write(typeLogFile, str string) bool
}

// ************** запрос ***************
type Requester interface {
	CommonTransmitter
	GetData() []byte
	SetData([]byte)
}

//************** ответ ***************

type Responser interface {
	CommonTransmitter
	GetData() any
	SetData(any)
	GetTaskId() string
	SetTaskId(string)
	GetSource() string
	SetSource(string)
	GetError() error
	SetError(error)
}

type CommonTransmitter interface {
	GetId() string
	SetId(string)
}
