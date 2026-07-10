# FastLearn MVP - Payment Strategy Update

## Phase 1: MVP (Current - Telegram based)

### Why Telegram?

✅ **Fast to implement** - no integration needed
✅ **Works right now** - Kaspi/PayBox require verification
✅ **Validates demand** - see who actually wants to pay
✅ **No compliance issues** - manual process, safer for MVP
✅ **Easy to scale** - Phase 2 can add automated payments

---

## How It Works

### User Flow

1. **User hits limit** (3 videos/month on free tier)
   ```
   Error: "Monthly limit reached (3/3)"
   "Want more? Contact us: t.me/fastlearn_support"
   ```

2. **User clicks Telegram link** 
   ```
   Opens: https://t.me/fastlearn_support
   Pre-filled: "Hi! I want to upgrade my FastLearn plan"
   ```

3. **User sends message**
   ```
   "Hi! I want Pro plan"
   ```

4. **You respond in Telegram**
   ```
   "Great! Pro plan is $2.99/month. 
    Send me your email and I'll activate it."
   ```

5. **You update database** (one SQL query)
   ```sql
   UPDATE subscriptions 
   SET plan='pro', status='active'
   WHERE user_id='their_user_id'
   ```

6. **User can upload 30 videos/month** ✅

---

## Implementation Details

### Backend Changes (Already Done ✅)

**File:** `backend/services/video.go`
```go
// MVP: All users have free tier limit of 3 videos/month
freeLimit := 3

if usage.VideosCount >= freeLimit {
    return false, "Monthly limit reached. Contact us via Telegram", nil
}
```

**File:** `backend/api/handler.go`
```go
// New endpoint for upgrade info
func (h *Handler) GetUpgradeInfo(c echo.Context) error {
    return c.JSON(200, map[string]interface{}{
        "contact": map[string]string{
            "type": "telegram",
            "username": "@fastlearn_support",  // ← CHANGE THIS
            "url": "https://t.me/fastlearn_support",  // ← CHANGE THIS
        },
        "plans": []map[string]string{
            {"name": "Pro", "price": "$2.99/month", "videos": "30"},
            {"name": "Premium", "price": "$7.99/month", "videos": "Unlimited"},
        },
    })
}
```

### Frontend Changes (Already Done ✅)

**File:** `frontend/src/views/HomePage.vue`
```vue
<div v-if="usage?.videos_count >= 3" class="upgrade-notice">
    <h3>🎯 Monthly limit reached</h3>
    <a :href="telegramLink" target="_blank" class="btn-upgrade-telegram">
        💬 Upgrade via Telegram
    </a>
</div>
```

---

## Configuration

### Step 1: Create Telegram Bot (optional)
If you want a bot to automate responses later:

1. Message @BotFather in Telegram
2. Create new bot: `/newbot`
3. Name: FastLearn Support
4. Username: `@fastlearn_support_bot`
5. Copy token: `123456:ABC-DEF1234ghIkl-zyx57W2v1u123ew11`

### Step 2: Update Code

**In `backend/api/handler.go` and `frontend/src/views/HomePage.vue`:**

Change:
```
@fastlearn_support  → @YOUR_USERNAME
https://t.me/fastlearn_support → https://t.me/YOUR_USERNAME
```

### Step 3: Manual Process (MVP)

Keep notifications in Telegram:
1. Telegram personal account: @fastlearn_support
2. Join FastLearn Telegram group / channel for notifications
3. When user messages, respond and update DB

---

## Database Updates (Manual)

### Pro Plan Activation

```sql
UPDATE subscriptions 
SET plan='pro', 
    status='active',
    current_period_start=CURRENT_TIMESTAMP,
    current_period_end=CURRENT_TIMESTAMP + INTERVAL '30 days'
WHERE user_id='user_id_here';
```

### Premium Plan Activation

```sql
UPDATE subscriptions 
SET plan='premium', 
    status='active',
    current_period_start=CURRENT_TIMESTAMP,
    current_period_end=CURRENT_TIMESTAMP + INTERVAL '30 days'
WHERE user_id='user_id_here';
```

### Cancel Subscription

```sql
UPDATE subscriptions 
SET plan='free', status='active'
WHERE user_id='user_id_here';
```

---

## Pricing Tiers (MVP)

| Plan | Price | Videos/mo | Duration | Support |
|------|-------|-----------|----------|---------|
| **Free** | $0 | 3 | 30 min max | Telegram |
| **Pro** | $2.99/mo | 30 | 3 hours max | Telegram |
| **Premium** | $7.99/mo | Unlimited | Unlimited | Telegram |

---

## Limitations & Next Steps

### MVP Limitations
❌ No automatic billing
❌ No payment history in app
❌ No invoice generation
❌ Manual activation only

### Phase 2 (Automated Payments)

Once you have paying customers:
1. Integrate PayBox API
2. Add webhooks for payment callbacks
3. Automatic subscription activation
4. Invoice/receipt generation
5. Recurring billing

---

## Conversation Example

```
User: Hi! I want to upgrade my FastLearn plan

You: Great! Which plan would you like?
   - Pro: 30 videos/month for $2.99/month
   - Premium: Unlimited for $7.99/month

User: I want Pro

You: Perfect! Please send me your Kaspi/PayBox details 
(or payment screenshot once they're set up in Phase 2)

You: (Update DB with SQL)
UPDATE subscriptions SET plan='pro' WHERE user_id='...';

You: Done! Your Pro plan is now active. 
You can upload 30 videos this month! 🎉

User: Thanks!
```

---

## Tracking Payments (MVP)

Create a simple spreadsheet:

| Date | Username | Email | Plan | Price | Notes |
|------|----------|-------|------|-------|-------|
| 2024-06-29 | user123 | email@test.com | Pro | $2.99 | Paid via Telegram |
| 2024-06-30 | user456 | email2@test.com | Premium | $7.99 | Paid via Telegram |

---

## Why This Works for MVP

1. **Validates demand** - You see who actually pays
2. **No risk** - You don't handle payments directly
3. **Compliant** - No payment processing regulations
4. **Fast** - No integration with Kaspi/PayBox
5. **Flexible** - Easy to change pricing/plans
6. **Personal** - Builds relationship with early users
7. **Safe** - You can refund if needed

---

## Transition to Phase 2

When you're ready for automated payments:

1. Kaspi / PayBox API integration
2. Webhook handling for payment confirmations
3. Automatic subscription activation
4. Billing history in app
5. Recurring billing setup

But for now, **keep it simple and manual**. You'll learn what users actually want.

---

**Remember:** The best feature is one that's shipped and used. This MVP approach gets you live faster! 🚀
