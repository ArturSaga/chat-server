package service

//go:generate sh -c "rm -rf mocks && mkdir -p mocks"
//go:generate /Users/artursagataev/GolandProjects/chat-server/bin/minimock -i ChatService -o ./mocks/ -s "_minimock.go"
//go:generate /Users/artursagataev/GolandProjects/chat-server/bin/minimock -i MessageService -o ./mocks/ -s "_minimock.go"
