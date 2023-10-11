package service

import (
	"context"
	"log"
	"strconv"
	"time"

	"github.com/azinudinachzab/scr-syky-tech-test/model"
	"gopkg.in/gomail.v2"
)

func (s *AppService) SendDailyBirthdayPromo(ctx context.Context, date time.Time) {
	users, err := s.fetchUsers(ctx, date.Format("2006-01-02"))
	if err != nil {
		log.Println("SendDailyBirthdayPromo fetchUsers: ", err)
		return
	}
	promos, err := s.generatePromoCode(ctx, users, date)
	if err != nil {
		log.Println("SendDailyBirthdayPromo generatePromoCode: ", err)
		return
	}
	s.SendNotification(ctx, users, promos)
	return
}

func (s *AppService) fetchUsers(ctx context.Context, date string) ([]model.User, error) {
	users, err := s.repo.GetUsersByFilter(ctx, map[string]string{
		"birth_date": date,
		"verification_status": strconv.Itoa(model.VerificationStatusVerified),
	})
	if err != nil {
		return nil, err
	}

	return users, nil
}

func (s *AppService) generatePromoCode(ctx context.Context, usr []model.User, date time.Time) ([]model.Promo, error) {
	promos := make([]model.Promo, 0)
	for _, val := range usr {
		promos = append(promos, model.Promo{
			UserID:     val.ID,
			PromoCode:  strconv.FormatUint(uint64(time.Now().UnixNano()), 10),
			Status:     model.PromoStatusActive,
			ExpiryAt:   date.Add(24 * time.Hour),
			Percentage: model.PromoPercentage,
		})
	}
	err := s.repo.StoreBulkPromo(ctx, promos)
	if err != nil {
		return nil, err
	}


	return promos, nil
}

func (s *AppService) SendNotification(ctx context.Context, u []model.User, p []model.Promo) {
	for i:=0;i<len(u);i++{
		go func(sender, pass string, usr model.User, prm model.Promo) {
			var (
				CONFIG_SMTP_HOST = "smtp.gmail.com"
				CONFIG_SMTP_PORT = 587
				CONFIG_SENDER_NAME = "PT. Sayakaya <"+sender+">"
				CONFIG_AUTH_EMAIL = sender
				CONFIG_AUTH_PASSWORD = pass
			)

			gender := "Bapak / Saudara"
			if usr.Gender == "P" {
				gender = "Ibu / Saudari"
			}

			mailer := gomail.NewMessage()
			mailer.SetHeader("From", CONFIG_SENDER_NAME)
			mailer.SetHeader("To", usr.Email)
			//mailer.SetAddressHeader("Cc", "tralalala@gmail.com", "Tra Lala La")
			mailer.SetHeader("Subject", "Selamat Ulang Tahun!")
			mailer.SetBody("text/html", "Kepada " + gender + " " + usr.FullName +
				", anda berhak mendapatkan promo untuk pembelian semua produk sayakaya sebanyak satu kali " +
				"berlaku hingga: " + prm.ExpiryAt.Format(time.RFC822) + ", kode promo: " + prm.PromoCode)
			//mailer.Attach("./sample.png")

			dialer := gomail.NewDialer(
				CONFIG_SMTP_HOST,
				CONFIG_SMTP_PORT,
				CONFIG_AUTH_EMAIL,
				CONFIG_AUTH_PASSWORD,
			)

			err := dialer.DialAndSend(mailer)
			if err != nil {
				log.Println("failed to send to: ", usr.Email, err)
			}
		}(s.conf.GMailUsername, s.conf.GMailPassword, u[i], p[i])
	}
}
