# UI Flows — Web

## Platform

Web

## Core Flows

### Login Flow

1. User navigates to any authenticated page
2. Auth middleware detects no valid JWT cookie → redirect to `/login`
3. User enters email and password
4. Submits form (standard POST, no HTMX)
5. Server validates credentials against bcrypt hash
   - **Success:** Set httpOnly JWT cookie (24h expiry) → redirect to `/dashboard`
   - **Error:** Re-render login page with "Invalid email or password" toast

### Create Webhook Flow

1. User is on Dashboard page
2. Clicks "New Webhook" button → inline form or modal appears (HTMX swap)
3. User fills in name (required) and description (optional)
4. Submits form (HTMX POST)
   - **Success:** New webhook card appears in the list (HTMX swap), toast "Webhook created"
   - **Error:** Validation error shown inline ("Name is required", "Name too long")
5. User can immediately copy the generated webhook URL via copy button

### View Webhook Requests Flow

1. User clicks a webhook card on Dashboard
2. Navigates to `/webhooks/{uuid}` — Webhook Detail page
3. Page shows webhook info (name, description, URL with copy button) + list of captured requests
4. SSE connection opens automatically via HTMX SSE extension
5. When a new request arrives externally:
   - Server pushes HTML fragment via SSE
   - HTMX prepends new request row to the list with fade-in animation (no page refresh)
6. User clicks a request row to navigate to full request detail

### View Request Detail Flow

1. User clicks a request row on Webhook Detail page
2. Navigates to `/webhooks/{uuid}/requests/{uuid}`
3. Full request displayed:
   - Method badge + full URL path
   - Headers table (key-value pairs)
   - Body content (pretty-printed if JSON/XML)
   - Query parameters table
   - Source IP and timestamp
4. Copy buttons available for: headers (as JSON), body (raw), full webhook URL

### Delete Webhook Flow

1. User clicks delete button (trash icon) on a webhook card
2. Confirmation modal appears: "Delete this webhook and all captured requests? This cannot be undone."
3. User confirms
   - **Success:** Webhook card removed from list (HTMX swap), toast "Webhook deleted"
   - **Cancel:** Modal closes, no action taken

### Edit Webhook Flow

1. User clicks edit button on webhook detail page
2. Name and description become editable (HTMX swap to form)
3. User modifies fields and submits
   - **Success:** Updated info displayed (HTMX swap), toast "Webhook updated"
   - **Error:** Validation errors shown inline

### Change Password Flow

1. User navigates to `/settings` via navbar link
2. Fills in current password, new password, confirm new password
3. Submits form (HTMX PUT)
   - **Success:** Toast "Password updated successfully"
   - **Error:** Inline errors ("Current password is incorrect", "Passwords don't match", "Minimum 8 characters")

### Toggle Theme Flow

1. User clicks sun/moon icon in navbar
2. Vanilla JS toggles `dark` class on `<html>` element
3. Preference saved to `localStorage`
4. On every page load, JS reads `localStorage` and applies theme before first render (prevents flash of wrong theme)

### Logout Flow

1. User clicks "Logout" in navbar
2. POST to `/logout` clears JWT cookie
3. Redirect to `/login`

## Screen States

| Screen | Empty State | Loading State | Error State |
|--------|-------------|---------------|-------------|
| Dashboard | "No webhooks yet. Create your first one!" with prominent CTA button | HTMX loading indicator (spinner on affected element) | Toast with error message |
| Webhook Detail | "No requests captured yet. Send a request to:" with webhook URL and copy button | HTMX loading indicator | Toast with error message |
| Request Detail | N/A (always has data if navigated to) | HTMX loading indicator | "Request not found" message with link back to webhook |
| Login | Always shows form | Submit button disabled with spinner | Inline "Invalid email or password" message |
| Settings | Always shows form | Submit button disabled with spinner | Inline validation errors |

## Key Interactions

| Interaction | Implementation |
|-------------|----------------|
| Copy to clipboard | Vanilla JS `navigator.clipboard.writeText()` → button text changes to "Copied!" for 2s, then reverts |
| SSE live updates | HTMX SSE extension (`hx-ext="sse"`), new rows prepend to request list with CSS fade-in animation |
| Delete confirmation | Custom modal component triggered by HTMX or `hx-confirm` attribute |
| Theme toggle | No page reload — JS toggles class + saves to localStorage |
| Method badge colors | Tailwind utility classes per method: GET=green, POST=blue, PUT=amber, PATCH=purple, DELETE=red |
| Form submissions | HTMX handles most forms — partial page swaps without full reload |
| Navigation | Standard server-side navigation (full page loads) — HTMX used for in-page actions only |
