package users

import (
	imageModel "admin/models/image"
	postModel "admin/models/post"
	"admin/routes"
	"fmt"

	"admin/routes/index"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

func RemoveImageController(c *fiber.Ctx) error {
	ImageID, err := strconv.ParseUint(c.Params("imageID"), 10, 64)

	if err != nil {
		return c.JSON(routes.HTTPResponse{
			Message: "Invalid Parameter Recieved",
			Success: false,
			Data:    nil,
		})
	}

	PostID, err := strconv.ParseUint(c.Params("postID"), 10, 64)

	if err != nil {
		return c.JSON(routes.HTTPResponse{
			Message: "Invalid Parameter Recieved",
			Success: false,
			Data:    nil,
		})
	}

	if err := imageModel.RemoveImage(ImageID, PostID); err != nil {
		return c.JSON(routes.HTTPResponse{
			Message: "Internal Server Error",
			Success: false,
			Data:    nil,
		})
	}

	return c.JSON(routes.HTTPResponse{
		Message: fmt.Sprintf("Deleted Image With ID %d", ImageID),
		Success: true,
		Data:    nil,
	})
}

func AddImageGetController(c *fiber.Ctx) error {
	user := index.GetSession(c).Get("User")
	postID, err := strconv.ParseUint(c.Params("postID"), 10, 64)

	if user == nil {
		return c.Redirect("/sign-in?err=please sign in")
	}

	if err != nil {
		return c.Redirect("/site/main?err=invalid parameter recieved")
	}

	post, err := postModel.GetPostById(postID)

	if err != nil {
		return c.Redirect("/site/main?err=post may have been deleted")
	}

	return c.Render("sites/main/image", fiber.Map{
		"Title": "Add Image",
		"User":  user,
		"Post":  post,
		"Breadcrumbs": []map[string]string{
			{"text": "Home", "linkTo": "/"},
			{"text": "Account", "linkTo": "/users/account"},
			{"text": "Main", "linkTo": "/site/main"},
			{"text": fmt.Sprintf("Post %d", postID), "linkTo": fmt.Sprintf("/site/main/post/%d", postID)},
			{"text": "Add Image", "linkTo": fmt.Sprintf("/site/main/post/%d/add-image", postID)},
		},
		"Error":   c.Query("err"),
		"Success": c.Query("s"),
	}, "layouts/main")
}

func AddImagePostController(c *fiber.Ctx) error {
	var body imageModel.Image

	postID, err := strconv.ParseUint(c.Params("postID"), 10, 64)

	if err != nil {
		return c.Redirect("/site/main?err=invalid parameter recieved")
	}

	if err := c.BodyParser(&body); err != nil {
		c.Redirect(fmt.Sprintf("/site/main/post/%d/add-image?err=unable to parse body", postID))
	}

	if err := imageModel.SaveNewImage(postID, body.URL); err != nil {
		c.Redirect(fmt.Sprintf("/site/main/post/%d/add-image?err=unable to save image", postID))
	}

	fmt.Println("Saved new image")
	return c.Redirect(fmt.Sprintf("/site/main/post/%d/add-image?s=image has been saved", postID))
}