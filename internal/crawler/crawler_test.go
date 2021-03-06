package crawler

import (
	"testing"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

func TestLoadConfig(t *testing.T) {

	entryPoint := viper.GetStringSlice("entryPoint")
	maxDepth := viper.GetInt("maxDepth")
	parallelism := viper.GetInt("parallelism")
	delay := viper.GetInt("delay")


	assert.True(t, len(entryPoint) > 0, "Should be greater than zero")
	assert.True(t, maxDepth > 0, "should be greater than zero")
	assert.True(t, parallelism > 0, "Should be greater than zero")
	assert.True(t, delay > 0, "Should be greater than zero")

}


func TestValidatePageUrl(t *testing.T) {

	crawler := New()

	url := "http://www.domain.com/en"
	assert.True(t, !crawler.isBlacklist(url), "This should be valid website url")


	php := "http://www.domain.com/en.php"
	assert.True(t, !crawler.isBlacklist(php), "This should be valid website url (.php)")

	asp := "http://www.domain.com/en.asp"
	assert.True(t, !crawler.isBlacklist(asp), "This should be valid website url (.asp)")


	aspx := "http://www.domain.com/en.aspx"
	assert.True(t, !crawler.isBlacklist(aspx), "This should be valid website url (.aspx)")

	jsp := "http://www.domain.com/en.jsp"
	assert.True(t, !crawler.isBlacklist(jsp), "This should be valid website url (.jsp)")


	html := "http://www.domain.com/en.html"
	assert.True(t, !crawler.isBlacklist(html), "This should be valid website url (.html)")


	htm := "http://www.domain.com/en.jsp"
	assert.True(t, !crawler.isBlacklist(htm), "This should be valid website url (.htm)")



}





func TestValidateFileUrl(t *testing.T) {

	crawler := New()

	url := "http://www.domain.com/file.pdf"
	assert.True(t, crawler.isBlacklist(url), "This should be url endpoint point to file")



	url = "http://www.domain.com/file.docx"
	assert.True(t, crawler.isBlacklist(url), "This should be url endpoint point to file")

}

