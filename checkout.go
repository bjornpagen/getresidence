package getresidence

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/efficientgo/core/errors"
	"github.com/stripe/stripe-go/v75"
	"github.com/stripe/stripe-go/v75/client"
	"github.com/uptrace/bunrouter"
)

func (s *server) postDubaiCheckout() bunrouter.HandlerFunc {
	return func(w http.ResponseWriter, r bunrouter.Request) error {
		session, err := getSession(r)
		if err != nil {
			return errors.Wrap(err, "get session")
		}

		idstr, err := s.branca.DecodeToString(session)
		if err != nil {
			return errors.Wrap(err, "decode session")
		}

		id, err := strconv.ParseInt(idstr, 10, 64)
		if err != nil {
			return errors.Wrap(err, "parse id")
		}

		name, email, phone, err := s.db.getOnboarding(id)
		if err != nil {
			return errors.Wrap(err, "get onboarding")
		}

		if name == "" || email == "" || phone == "" {
			return errors.New("name, email, or phone is empty")
		}

		var customer *stripe.Customer
		customer, err = customerBySession(s.stripe, session)
		if err != nil {
			customer, err = createCustomer(s.stripe, name, email, phone, session)
			if err != nil {
				return errors.Wrap(err, "create stripe customer")
			}
		} else {
			_, err = s.stripe.Customers.Update(customer.ID, &stripe.CustomerParams{
				Name:  stripe.String(name),
				Email: stripe.String(email),
				Phone: stripe.String(phone),
			})
			if err != nil {
				return errors.Wrap(err, "update stripe customer")
			}
		}

		const PRICE_ID = "price_1NfSN6Be2isi3fifFO60qs76"

		checkoutSession, err := s.stripe.CheckoutSessions.New(&stripe.CheckoutSessionParams{
			SuccessURL: stripe.String("https://getresidence.org/dubai/success"),
			CancelURL:  stripe.String("https://getresidence.org/dubai/cancel"),
			Customer:   stripe.String(customer.ID),
			Mode:       stripe.String(string(stripe.CheckoutSessionModePayment)),
			LineItems: []*stripe.CheckoutSessionLineItemParams{
				{
					Quantity: stripe.Int64(1),
					Price:    stripe.String(PRICE_ID),
				},
			},
		})
		if err != nil {
			return errors.Wrap(err, "create stripe checkout session")
		}

		w.Header().Set("HX-Redirect", checkoutSession.URL)
		return nil
	}
}

func customerBySession(sc *client.API, session string) (*stripe.Customer, error) {
	params := &stripe.CustomerSearchParams{}
	params.Query = fmt.Sprintf(`metadata['session']:'%s'`, session)
	params.Limit = stripe.Int64(1)

	iter := sc.Customers.Search(params)
	for iter.Next() {
		result := iter.Current()

		switch result, ok := result.(*stripe.Customer); {
		case !ok:
			return nil, errors.New("unexpected type")
		default:
			return result, nil
		}
	}

	if err := iter.Err(); err != nil {
		return nil, errors.Wrap(err, "search stripe customer")
	}

	return nil, errors.New("no customer found")
}

func createCustomer(sc *client.API, name, email, phone, session string) (*stripe.Customer, error) {
	params := &stripe.CustomerParams{
		Name:  stripe.String(name),
		Email: stripe.String(email),
		Phone: stripe.String(phone),
	}
	params.AddMetadata("session", session)

	customer, err := sc.Customers.New(params)
	if err != nil {
		return nil, errors.Wrap(err, "create stripe customer")
	}

	return customer, nil
}
