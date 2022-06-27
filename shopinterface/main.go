package main

import (
	// gin

	"log"
	"math/rand"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/joesjo/grpc-store/shopinterface/serviceclient"
)

const (
	PORT = ":5000"
)

type CreateItem struct {
	Name     string `json:"name" binding:"required"`
	Quantity int32  `json:"quantity"`
}

type StockItem struct {
	ItemId   string `json:"item_id" binding:"required"`
	Quantity int32  `json:"quantity" binding:"required"`
}

func main() {
	serviceclient.Init()
	router := gin.Default()
	router.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "Hello World!",
		})
	})
	router.GET("/inventory", func(c *gin.Context) {
		inventoryList, err := serviceclient.GetInventory()
		if err != nil {
			c.JSON(500, gin.H{
				"error": err.Error(),
			})
			return
		}
		c.JSON(200, inventoryList)
	})
	router.POST("/inventory/create", func(c *gin.Context) {
		var createItem CreateItem
		if err := c.BindJSON(&createItem); err != nil {
			c.JSON(500, gin.H{
				"error": err.Error(),
			})
			return
		}
		itemId, err := serviceclient.CreateItem(createItem.Name, createItem.Quantity)
		if err != nil {
			c.JSON(500, gin.H{
				"error": err.Error(),
			})
			return
		}
		c.JSON(200, gin.H{
			"itemId": itemId,
		})
	})
	router.POST("/inventory/stock", func(c *gin.Context) {
		var stockItem StockItem
		if err := c.BindJSON(&stockItem); err != nil {
			c.JSON(500, gin.H{
				"error": err.Error(),
			})
			return
		}
		err := serviceclient.StockItem(stockItem.ItemId, stockItem.Quantity)
		if err != nil {
			c.JSON(500, gin.H{
				"error": err.Error(),
			})
			return
		}
		c.JSON(200, gin.H{
			"message": "Stock item successful",
		})
	})
	router.POST("/inventory/purchase", func(c *gin.Context) {
		quantity, err := strconv.Atoi(c.PostForm("quantity"))
		if err != nil {
			c.JSON(500, gin.H{
				"error": err.Error(),
			})
			return
		}
		err = serviceclient.PurchaseItem(c.PostForm("itemId"), int32(quantity))
		if err != nil {
			c.JSON(500, gin.H{
				"error": err.Error(),
			})
			return
		}
		c.JSON(200, gin.H{
			"message": "Purchase item successful",
		})
	})
	router.GET("/stresstest", func(c *gin.Context) {
		stressTest()
		c.JSON(200, gin.H{
			"message": "Stress test started",
		})
	})
	router.Run(PORT)
}

func stressTest() {
	adjectives := []string{"autumn", "hidden", "bitter", "misty", "silent", "empty", "dry", "dark", "summer", "icy", "delicate", "quiet", "white", "cool", "spring", "winter", "patient", "twilight", "dawn", "crimson", "wispy", "weathered", "blue", "billowing", "broken", "cold", "damp", "falling", "frosty", "green", "long", "late", "lingering", "bold", "little", "morning", "muddy", "old", "red", "rough", "still", "small", "sparkling", "throbbing", "shy", "wandering", "withered", "wild", "black", "young", "holy", "solitary", "fragrant", "aged", "snowy", "proud", "floral", "restless", "divine", "polished", "ancient", "purple", "lively", "nameless"}
	nouns := []string{"waterfall", "river", "breeze", "moon", "rain", "wind", "sea", "morning", "snow", "lake", "sunset", "pine", "shadow", "leaf", "dawn", "glitter", "forest", "hill", "cloud", "meadow", "sun", "glade", "bird", "brook", "butterfly", "bush", "dew", "dust", "field", "fire", "flower", "firefly", "feather", "grass", "haze", "mountain", "night", "pond", "darkness", "snowflake", "silence", "sound", "sky", "shape", "surf", "thunder", "violet", "water", "wildflower", "wave", "water", "resonance", "sun", "wood", "dream", "cherry", "tree", "fog", "frost", "voice", "paper", "frog", "smoke", "star"}

	for i := 0; i < 100; i++ {
		go func() {
			// create item with random name and quantity
			name := adjectives[rand.Intn(len(adjectives))] + " " + nouns[rand.Intn(len(nouns))]
			quantity := rand.Int31n(100)
			itemId, err := serviceclient.CreateItem(name, quantity)
			if err != nil {
				log.Println(err)
				return
			}
			log.Println("Created item: " + itemId)
		}()
	}
}
