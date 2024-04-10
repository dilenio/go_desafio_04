package limiter

import (
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/dilenio/desafio04/cmd/configs"
	"github.com/dilenio/desafio04/storage"
	"github.com/go-redis/redis"
)

type RateLimiter struct {
	redisClient *storage.RedisStorage
}

func NewRateLimiter() *RateLimiter {
	config, err := configs.LoadConfig(".")
	if err != nil {
		panic(err)
	}

	redisClient, err := storage.NewRedisStorage(config.RedisHost, "", 0)
	if err != nil {
		log.Fatalf("Erro ao conectar com o Redis: %v", err)
	}

	return &RateLimiter{
		redisClient: redisClient,
	}
}

func (rl *RateLimiter) processRequest(w http.ResponseWriter, key, keyType string) bool {
	if rl.isBlocked(key) {
		http.Error(w, "you have reached the maximum number of requests or actions allowed within a certain time frame", http.StatusTooManyRequests)
		return false
	}
	if rl.checkRateLimit(key, keyType) {
		rl.block(key, keyType)
		http.Error(w, "you have reached the maximum number of requests or actions allowed within a certain time frame", http.StatusTooManyRequests)
		return false
	}
	return true
}

func (rl *RateLimiter) LimitHandler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		config, _ := configs.LoadConfig(".")

		apiKey := r.Header.Get("API_KEY")
		ip := strings.Split(r.RemoteAddr, ":")[0]

		var key, keyType string
		if apiKey == config.TokenAllowed {
			key = "token:" + apiKey
			keyType = "token"
		} else {
			key = "ip:" + ip
			keyType = "ip"
		}

		if !rl.processRequest(w, key, keyType) {
			return // Se a função retornar false, a requisição é interrompida
		}

		next.ServeHTTP(w, r)
	})
}

func (rl *RateLimiter) checkRateLimit(key string, tokenOrIp string) bool {
	config, err := configs.LoadConfig(".")
	if err != nil {
		panic(err)
	}

	val, err := rl.redisClient.Get(key)

	if err == redis.Nil {
		if tokenOrIp == "ip" {
			rl.redisClient.Set(key, "1", time.Duration(config.RequestsByIp)*time.Second)
		} else {
			rl.redisClient.Set(key, "1", time.Duration(config.RequestsByToken)*time.Second)
		}
		return false
	}

	count, err := strconv.Atoi(val)
	if err != nil {
		log.Printf("Erro ao converter o valor do contador: %v\n", err)
		return false
	}

	var requests int

	if tokenOrIp == "ip" {
		requests = config.RequestsByIp
	} else {
		requests = config.RequestsByToken
	}

	if count >= requests {
		return true
	}

	rl.redisClient.Incr(key)
	return false
}

func (rl *RateLimiter) block(key string, tokenOrIp string) {
	config, err := configs.LoadConfig(".")
	if err != nil {
		panic(err)
	}
	var timeBlocked int
	if tokenOrIp == "ip" {
		timeBlocked = config.TimeBlockedByIp
	} else {
		timeBlocked = config.TimeBlockedByToken
	}
	rl.redisClient.Set(key+":blocked", "1", time.Duration(timeBlocked)*time.Second)
}

func (rl *RateLimiter) isBlocked(key string) bool {
	_, err := rl.redisClient.Get(key + ":blocked")
	return err == nil
}
