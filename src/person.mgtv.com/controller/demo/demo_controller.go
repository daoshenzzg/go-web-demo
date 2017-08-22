package controller

import (
	"person.mgtv.com/framework/mvc"
	"person.mgtv.com/framework/resultcode"
	demoServcie "person.mgtv.com/service/demo"
	"person.mgtv.com/framework/logs"
)

type DemoController struct {
	mvc.Controller
	demoServcie *demoServcie.DemoService
}

func (c *DemoController) InitController() {
	c.demoServcie = demoServcie.NewDemoService()
}

func (c *DemoController) URLMapping() {
	c.Mapping("get_feed", c.getFeed)
	c.Mapping("get_key", c.getKey)
	c.Mapping("is_followed", c.isFollowed)
	c.Mapping("multi_commit", c.multiCommit)
}

// test mysql get
func (c *DemoController) getFeed() {
	feedId := c.Get("feedId")
	if feedId == "" {
		c.Error(resultcode.ERROR_1001)
		return
	}

	feed, err := c.demoServcie.GetFeed(feedId)
	if err != nil {
		c.Error(resultcode.ERROR, err)
		return
	}
	c.Data("feed", feed)
}

// test redis get
func (c *DemoController) getKey() {

	value, err := c.demoServcie.GetKey("key")
	if err != nil {
		c.Error(resultcode.ERROR, err)
		return
	}
	c.Data("value", value)
}

// test http get
func (c *DemoController) isFollowed() {
	uid := c.Get("uid")
	if uid == "" {
		c.Error(resultcode.ERROR_1001)
		return
	}

	artistId := c.Get("artistId")
	if artistId == "" {
		c.Error(resultcode.ERROR_1001)
		return
	}

	isFollowed, _ := c.demoServcie.IsFollowed(uid, artistId)
	c.Data("isFollowed", isFollowed)
}

// test mysql transaction
func (c *DemoController) multiCommit() {
	err := c.demoServcie.MutiCommit()
	if err != nil {
		logs.GetLogger("system").Errorf("Multi commit error", err)
	}
	c.Success()
}
