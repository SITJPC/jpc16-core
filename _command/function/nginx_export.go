package function

import (
	"fmt"
	"os"
	"strings"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"jpc16-core/common/mng"
	"jpc16-core/type/collection"
	"jpc16-core/util/log"
)

func ConfigExport() {
	// * Find groups
	var groups []*collection.Group
	if err := mng.GroupCollection.SimpleFind(&groups, bson.M{}); err != nil {
		log.Fatal("Unable to get groups", err)
	}

	// * Find players
	var players []*collection.Player
	if err := mng.PlayerCollection.SimpleFind(&players, bson.M{}); err != nil {
		log.Fatal("Unable to get players", err)
	}

	// * Construct group map
	groupMap := make(map[primitive.ObjectID]*collection.Group)
	for _, group := range groups {
		groupMap[*group.ID] = group
	}

	nginxFile, err := os.Create("local/nginx.conf")
	if err != nil {
		log.Error("Failed to create file", err)
	}
	defer nginxFile.Close()

	sftpFile, err := os.Create("local/sftp.conf")
	if err != nil {
		log.Error("Failed to create file", err)
	}

	template := `
server {
	listen 80;
	server_name %s-%s.sjpc.me;
	root /var/www/data/%s-%s;
}
`
	// * Append to string
	var nginxData string
	var sftpData string
	for _, player := range players {
		// * Generate group code
		splittedGroupName := strings.Split(*groupMap[*player.GroupId].Name, " ")
		var groupCode string
		for _, word := range splittedGroupName {
			groupCode += string(word[0])
		}
		groupCode = strings.ToLower(groupCode)

		// * Append to nginx data
		nginxData += fmt.Sprintf(
			template,
			strings.ToLower(*player.Nickname),
			groupCode,
			strings.ToLower(*player.Nickname),
			groupCode,
		)

		// * Append to sftp data
		sftpData += fmt.Sprintf("%s-%s:%s::100\n", strings.ToLower(*player.Nickname), groupCode, *player.Pin)
	}

	// * Write to file
	if _, err := nginxFile.Write([]byte(nginxData)); err != nil {
		log.Debug("Unable to export nginx config", "file name", nginxFile.Name())
	}
	if _, err := sftpFile.Write([]byte(sftpData)); err != nil {
		log.Debug("Unable to export sftp config", "file name", sftpFile.Name())
	}

	log.Debug("Successfully exported", "nginx file name", nginxFile.Name(), "sftp file name", sftpFile.Name())
}
