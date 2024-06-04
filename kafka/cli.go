package kafka

import (
	"log"
)

func Cli(aliasName ...string) MessageProducer {
	name := "default"
	if aliasName != nil && len(aliasName) > 0 && aliasName[0] != "" {
		name = aliasName[0]
	}
	p, ok := producers.Load(name)
	if !ok {
		log.Panicf("no <%s> kafka client found(need NewProducer first)\n", name)
		return nil
	}
	return p.(*Producer)
}

func CliV2(aliasName ...string) MessageProducer {
	name := "default"
	if aliasName != nil && len(aliasName) > 0 && aliasName[0] != "" {
		name = aliasName[0]
	}
	p, ok := v2Producers.Load(name)
	if !ok {
		log.Panicf("no <%s> kafka client found(please invoke CreateProducer first)\n", name)
		return nil
	}
	return p.(*messageProducer)
}
