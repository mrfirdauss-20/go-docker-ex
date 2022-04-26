type Storage struct {
	redisClient *redis.Client
}

type Config struct {
	redisClient *redis.Client `validate:"nonnil"`
}

func New(cfg Config) (*Storage, error) {
	err:= validator.Validate(cfg)
	if err!=nil{
		return nil, err
	}
	s:= &Storage{redisClient: cfg.redisClient}
	return s,nil
}

func (s *Storage) GetRandomQuestion(ctx context.Context) (*core.Question,error){

	//res, err := rdb.Do(ctx, "set", "key", "value").Result()
	//val2, err := rdb.Get(ctx, "key2").Result()
	var questions []core.Question
	i int :=0
	atr,err := s.redisClient.DoContext(ctx,"get","questions:0")
	for atr!= nil{
		CorIdx = s.redisClient.DoContext(ctx,"get","questions:"+i+":correct_index")
		Ans = s.redisClient.DoContext(ctx,"get","questions:"+i+":answers")
		var q core.Question ={
			Problem: atr.Result(),
			CorrectIndex: CorIdx.Result(),
			Answers: Ans.Result()
		}
		questions = append(questions,q)
		i++;
		atr,err= s.redisClient.DoContext(ctx,"get","questions:"+i)
		if(err!=nil){
			continue
		}
	}

	if len(questions) ==0{
		return nil, nil
	}
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	idx := r.Intn(len(questions))
	return &questions[idx],nil
}