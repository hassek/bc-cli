package templates

const ActiveSubscriptionsTemplate = `=== Your Active Subscriptions ===

{{range .Subscriptions}}{{if eq .Status "active"}}┌─ {{.Tier | upper}} ✓
│  Status: {{.Status | upper}}
{{if .StartedAt}}│  Started: {{.StartedAt}}
{{end}}
└─ ID: {{.ID}}

{{end}}{{end}}{{if .HasActive}}✓ = Active subscription

{{end}}{{repeat "=" 60}}

`

const SubscriptionDetailsTemplate = `
{{repeat "=" 60}}
{{.Name}}
{{repeat "=" 60}}

Price: {{.Currency}} {{.Price}}/{{.BillingPeriod}}
Description: {{.Description}}
{{if .ActiveSub.ID}}
Status: {{.ActiveSub.Status | upper}}{{if eq .ActiveSub.Status "active"}} ✓{{end}}
{{if .ActiveSub.StartedAt}}Started: {{.ActiveSub.StartedAt}}
{{end}}
{{end}}
Features:
{{range .Features}}  • {{.}}
{{end}}
`

const OrderConfigIntroTemplate = `
{{repeat "─" 60}}
Let's configure your coffee order!
{{repeat "─" 60}}

How much coffee would you like per month?
You can order anywhere from {{.MinQuantity}} kg to {{.MaxQuantity}} kg.
`

const OrderSplitIntroTemplate = `{{repeat "─" 60}}

Would you like your coffee prepared different ways?
For example, you could get:
  • 2 kg whole bean + 3 kg ground for espresso
  • 2 kg ground for moka + 2 kg ground for v60 + 1 kg whole bean

Or keep it simple with everything the same way.`

const UniformOrderIntroTemplate = `{{repeat "─" 60}}

Great! Let's prepare all {{.TotalQuantity}} kg the same way.

`

const SplitOrderIntroTemplate = `{{repeat "─" 60}}

Great! Now let's split your {{.TotalQuantity}} kg into different
grinding preferences. You can have:
  • Whole beans (you grind at home)
  • Pre-ground for specific brewing methods
We'll help you allocate all {{.TotalQuantity}} kg across your preferences.`

const PreferenceHeaderTemplate = `{{repeat "─" 60}}
┌─ Preference #{{.PreferenceNum}} ──────────────────────────────────────────┐
│ {{printf "%-58s" (printf "Allocating from: %d kg total" .TotalQuantity)}} │{{if .LowRemaining}}
│ {{printf "%-58s" (printf "Remaining: %d kg ⚠️  (almost done!)" .Remaining)}} │{{else}}
│ {{printf "%-58s" (printf "Remaining: %d kg" .Remaining)}} │{{end}}
└────────────────────────────────────────────────────────────┘
`

const ProgressBarTemplate = `
┌────────────────────────────────────────────────────────────┐{{if ge .Current .Total}}
│ {{printf "%-58s" (printf "Progress: %s %d/%d kg ✓" (progressBar .Current .Total 30) .Current .Total)}} │{{else}}
│ {{printf "%-58s" (printf "Progress: %s %d/%d kg" (progressBar .Current .Total 30) .Current .Total)}} │{{end}}
└────────────────────────────────────────────────────────────┘`

const OrderSummaryTemplate = `Your Order Summary:
┌─────────────────────────────────────────────────────────┐
│ {{printf "%-55s" (printf "Tier: %s" .TierName)}} │
│ {{printf "%-55s" (printf "Total: %d kg/month" .TotalQuantity)}} │
│ {{printf "%-55s" (printf "Price: %s %.2f/%s" .Currency .TotalPrice .BillingPeriod)}} │
│ {{printf "%-55s" ""}} │
│ {{printf "%-55s" "How your coffee will be prepared:"}} │
{{range $i, $item := .LineItems}}│ {{printf "%-55s" (printf "   %d. %s" (add $i 1) $item)}} │
{{end}}└─────────────────────────────────────────────────────────┘
`

const CheckoutHeaderTemplate = `
{{repeat "─" 60}}
Opening checkout...
`

const SuccessMessageTemplate = `
🎉 Congratulations! Your subscription is now active!

📦 Your first shipment of {{.TotalQuantity}} kg of fresh {{.TierName}} coffee
   will be shipped within the next 7 days.

☕ Get ready for an amazing coffee experience!
`

const SuccessArtTemplate = `
MMMMMMMMMMMMMWXOdc;;;cOWMMMMMMMMMMMMMMMM
MMMMMMMMMMMXxc,...''..'xWMMMMMMMMMMMMMMM
MMMMMMMMMMXc.......,,'.'xNX0OKWMMMMMMMMM
MMMMMMMMMMNo.......;cc:''::,,;kWMMMMMMMM
MMMMMMMMMMMXl..';;:cc:,'',;,,oKMMMMMMMMM
MMMMMMMMMMMW0;.,,'.''';:cdxdlxNMMMMMMMMM
MMMMMMMMMWKo;...';clodxxdxxoc:dKX0O0NMMM
MMMMMMMMMWd....:okO000Oc':dloxdxxxdl0MMM
MWWMMMMMMMNOxxlokO00000OxkxccooxxddkXMMM
XolKWMMMMMMMWKc;oxO00KK0KKOdc:dO00kd0WMM
d..,oOKNNNXOo,...';coddlcdxl,oNMNOdokXWM
c.....'::;'..',......,:...'..cXXo:llccdK
l......;,.....;:;,.....':dxl,oNK:.....:d
k'.....;'.......',;;,..;0N0Kxd0x:'..',;k
Nd......;..........;o:..xXO0l''.,cooodxK
MNd'....,,.........;l;..dNKd....'xXNNWMM
MMWO:....'''.....',,'...okl....,xNMMMMMM
MMMMNkc'........''.....'cooolokXWMMMMMMM
MMMMMMWKko:,.......,;;cx0WMMMMMMMMMMMMMM
MMMMMMMMMMNKOxdddxk0KXWMMMMMMMMMMMMMMMMM
`

// Subscription Management Templates

const ManageNotAuthenticatedTemplate = `You must be logged in to manage subscriptions.

Please run: bc-cli login
`

const NoSubscriptionsTemplate = `You don't have any subscriptions yet.

To subscribe, run: bc-cli subscriptions
`

const NoActionsAvailableTemplate = `No actions available for this subscription.
`

const ManageSubscriptionHeaderTemplate = `
{{repeat "=" 60}}
Managing Subscription: {{.Tier | upper}}
{{repeat "=" 60}}

{{.StatusIcon}} Status: {{.Status | upper}}
{{if .StartedAt}}Started: {{.StartedAt}}
{{end}}
{{if .HasNextShipment}}Next Shipment: {{.NextShipment}}
{{end}}
{{if .HasPricing}}
Billing: {{.Price}} {{.Currency}}/{{.BillingPeriod}}
{{end}}
{{if .HasOrderDetails}}
Current Order Configuration:
  Total: {{.TotalQuantity}} kg per month
{{range $i, $item := .LineItems}}  {{add $i 1}}. {{$item}}
{{end}}{{end}}
`

const PauseWarningTemplate = `
⚠  Pausing your subscription will:
  • Stop upcoming shipments
  • Pause billing
  • Keep your preferences saved
  • You can resume anytime

`

const PauseConfirmWithDateTemplate = `
✓ Your subscription will be paused for {{.Months}} month(s)
  and automatically resume on {{.ResumeDate}}

`

const SubscriptionPausedTemplate = `
✓ Subscription paused successfully!
{{if .HasResumeDate}}
📅 Your subscription will automatically resume on {{.ResumeDate}}
{{else}}
💤 Your subscription is paused indefinitely. Use 'bc-cli manage' to resume.
{{end}}
`

const ResumeInfoTemplate = `
✓ Resuming your subscription will:
  • Restart shipments
  • Resume billing

`

const SubscriptionResumedTemplate = `
✓ Subscription resumed successfully!

📦 Your next shipment will be scheduled soon.
`

const UpdateSubscriptionHeaderTemplate = `
{{repeat "─" 60}}
Update Subscription Preferences
{{repeat "─" 60}}

`

const UpdatePreferencesSummaryTemplate = `
{{repeat "─" 60}}
New Subscription Preferences:
{{repeat "─" 60}}

Total: {{.TotalQuantity}} kg per month

How your coffee will be prepared:
{{range $i, $item := .LineItems}}  {{add $i 1}}. {{$item}}
{{end}}
{{repeat "─" 60}}

`

const SubscriptionUpdatedTemplate = `
✓ Subscription updated successfully!

📦 Your changes will take effect with your next shipment.
`

const CancelWarningTemplate = `
⚠  Warning: Cancelling your subscription will:
  • Stop all future shipments
  • End your billing cycle
  • Remove access to subscription benefits
  • This action cannot be easily undone

💡 Did you know? You can pause your subscription instead!
   Pausing keeps your preferences and lets you resume anytime.

`

const CancelDoubleConfirmTemplate = `
Please confirm once more that you want to cancel permanently.
`

const SubscriptionCancelledTemplate = `
✓ Subscription cancelled.

We're sorry to see you go! If you change your mind,
you can always start a new subscription with: bc-cli subscriptions
`

const ActionCancelledTemplate = `
{{.Action}} cancelled.
`
