package repository

//go:generate sh -c "rm -rf mocks && mkdir -p mocks"
//go:generate /Users/artursagataev/GolandProjects/chat-server/bin/minimock -i ChatRepository -o ./mocks/ -s "_minimock.go"
//go:generate /Users/artursagataev/GolandProjects/chat-server/bin/minimock -i MessageRepository -o ./mocks/ -s "_minimock.go"
