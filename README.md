### Лабораторная 2

1. Запустить postgres
```zsh
    docker-compose up -d
```

2. Сбилдить бинарь 
```zsh
    go build ./main.go              
```

3. Запустить программу
```zsh
    ./main --s ./senses.xml --sw ./service_words.xml
```

для примера, введите: "как торжествовать победу"


