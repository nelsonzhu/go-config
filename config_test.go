package config

import "testing"

type testTConfig struct {
	IsTLS        bool
	Addr         string
	ReadTimeOut  int
	WriteTimeOut int
	Pattern1     string `xml:"Patterns>Pattern1"`
	Pattern2     string `xml:"Patterns>Pattern2"`
	DBDSN        string `xml:"DB>DSN"`
	DBUID        string `xml:"DB>UID"`
	DBPassWord   string `xml:"DB>PassWord"`
	ConnPoolSize int    `xml:"DB>ConnPoolSize"`
}

var test_item = testTConfig{true, "test", 10, 10, "p1", "p2", "localhost:3306", "defualt", "", 400}

var testConfigDates = []interface{}{
	test_item,
	uint16(61374),
	int8(-54),
	uint8(254),
	string(`
	<testConfig>
	    <IsTLS>true</IsTLS>
	    <ReadTimeOut>10</ReadTimeOut>
	    <WriteTimeOut>10</WriteTimeOut>
	    <Patterns>
	        <Pattern1>/</Pattern1>
	        <Pattern2>/view</Pattern2>
	        <Pattern3>/show</Pattern3>
	    </Patterns>
	    <Apps>
	        <App AppID="test" Version="1.0" Expires="2016-01-01T00:00:00.00+08:00">
	            <Token>764427669080</Token>
	        </App>
	        <AppT AppID="test1" Version="1.0" Expires="2016-01-01T00:00:00.00+08:00">
	            <Token>6902cbd0d24</Token>
	        </AppToken>
	    </Apps>
	    <DB>
	        <DSN>localhost:3306</DSN>
	        <UID>default</UID>
	        <PassWord></PassWord>
	        <DBName>test</DBName>
	        <ConnPoolSize>400</ConnPoolSize>
	    </DB>
	</testConfig>`),
}

func Test_GeterAndSetter(t *testing.T) {
	for i, data := range testConfigDates {
		conf := NewConfig(data)
		v := conf.ConfValue()
		if data != v {
			t.Errorf("get and set error, case:%d set:%#v get:%#v", i, data, v)
		}
	}
}
