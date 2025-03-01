package entity

type UserAggregate struct {
	User               User
	Detail             *UserDetail
	AccountLine        *UserAccountLine
	StripeSubscription *SubscribeStripe
	StateCount         *StateCount
}
