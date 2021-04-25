package main

import (
	"context"
	"fmt"
	"os"

	"github.com/cvcio/elections-api/models/user"
	"github.com/cvcio/elections-api/pkg/config"
	"github.com/cvcio/elections-api/pkg/db"
	"github.com/cvcio/elections-api/pkg/mailer"
	"github.com/kelseyhightower/envconfig"
	log "github.com/sirupsen/logrus"
)

var (
	cfg                 = config.New()
	msgAccountActivated = `
<div style="font-family:'Segoe UI',Arial,Sans-Serif;font-size:10pt;">
    <p>
        Hi %s,
    </p>
    <p>
        Ο λογαριασμός σας επιβεβαιώθηκε!
    </p>
    <p>
		Για να συνδεθείτε -και κάθε φορά που θα θέλετε να συνδεθείτε- θα πρέπει να κάνετε authorize την εφαρμογή με τον λογαριασμό σας στο Twitter.
		Για να ολογκληρώσετε την εγγραφή σας ακολοθήστε το παρακάτω link: <a href="https://elections.mediawatch.io/?provider=email&method=invite&email=%s">Ολοκλήρωση Εγγραφής</a>.
    </p>
    <p>
        Για οποιαδήποτε διευκρύνηση μπορείτε να επικοινωνήσετε μαζί μας στο email <a href="mailto:press@mediawatch.io">press@mediawatch.io</a> ή στο +30 211 103 5100.
    </p>
    <p>
        Καλή συνέχεια,<br />
        η ομάδα του MediaWatch
    </p>
</div>
`
)

func activate(screenName string) {
	log.Infof("Activating user: %s\n", screenName)

	// ============================================== ==============
	// Start Mongo
	log.Debug("Initialize Mongo")
	dbConn, err := db.New(cfg.MongoURL(), cfg.Mongo.Path, cfg.Mongo.DialTimeout)
	if err != nil {
		log.Fatalf("Register DB: %v", err)
	}
	log.Debug("Connected to Mongo")
	defer dbConn.Close()

	// Create mail service"github.com/sirupsen/logrus"
	mail := mailer.New(
		cfg.SMTP.Server,
		cfg.SMTP.Port,
		cfg.SMTP.User,
		cfg.SMTP.Pass,
		cfg.SMTP.From,
		cfg.SMTP.FromName,
		cfg.SMTP.Reply,
	)

	u, err := user.ByScreenNanme(dbConn, screenName)
	if err != nil {
		log.Fatalf("Can't activate user (%s): %s", screenName, err.Error())
	}

	if u.Status == "active" {
		log.Infof("User already activated, exiting...")
		os.Exit(0)
	}

	u.Status = "active"
	_, err = user.Update(dbConn, u.IDStr, u)
	if err != nil {
		log.Fatalf("Can't update user (%s): %s", screenName, err.Error())
	}

	log.Infof("Sending email to user %s (%s)", u.ScreenName, u.Email)

	mailer.Message(
		context.Background(),
		mail,
		u.Email,
		"Account Activated | MediaWatch / EU Elections 2019",
		fmt.Sprintf(msgAccountActivated, u.FirstName, u.Email),
	)

	log.Infof("User activated %s", u.ScreenName)
}

func main() {
	// Read config from env variables
	err := envconfig.Process("", cfg)
	if err != nil {
		log.Fatalf("Error loading config: %s", err.Error())
	}

	methods := os.Args[1:]
	if methods[0] == "activate" && len(methods) == 2 {
		activate(methods[1])
	}

	log.Info("Exit")
}
