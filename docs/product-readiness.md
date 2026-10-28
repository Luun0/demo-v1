# Product readiness: SEN Support public preview

**Reviewed:** 28 October 2026. This is a product-readiness checklist, not a launch approval or a privacy/legal opinion.

## Current product state

- A public evaluation preview is available at <https://sen-support-platform.onrender.com/>. Its Render web service is on the Free plan in Ohio and deploys from `main`; the configured Neon branch is `production` in Ohio (`us-east-2`). Free Render services sleep when idle.
- The product supports guardian, assigned-specialist, and administrator workflows for a child profile, sessions, visibility-controlled diary entries, and access history. Administrators do not have diary-text access.
- Public account creation is available, but a specialist needs administrator approval before assignment. The public site does not have email delivery configured for password-reset links.
- The deployment is an academic prototype. Use fictional data only. No real child or family data is approved for this environment.
- There is no product analytics or completed usability study. Access logs are security/audit records, not a product-funnel report.

## Gates before real child or family data

1. Confirm the data controller/operator, purpose, lawful basis, guardian authority, and age-appropriate notice with a qualified reviewer.
2. Confirm whether the current Render and Neon Ohio regions and all backup/subprocessor locations are legally acceptable for the intended data. The present configuration is not an approval.
3. Define retention, backup deletion, data-subject request handling, incident contact, and operational ownership.
4. Provide a user-facing privacy notice that accurately names the operator, data, recipients, location, retention, and support contact.
5. Configure password recovery through a verified email delivery provider using an HTTPS API. Do not enable reset-link logging on the public service.
6. Verify the actual production deployment and access controls with synthetic accounts before inviting participants.

## Next usability evaluation

Recruit guardians and specialists separately. Observe participants completing these tasks with fictional data only:

1. Explain what the service is for and what data must not be entered.
2. Register in the correct role and describe what happens after specialist registration.
3. Create a sample child profile, assign an approved specialist, and schedule a session.
4. Add a diary entry, choose its visibility, and find the access history.

Record task completion, mistakes, questions, and points where participants stop. Do not infer success from time on page alone. The sample size and findings must be reported honestly; the existing seven-person exploratory survey is not a usability study.

## Product measurement, only after privacy review

If analytics are later approved, use aggregate events without emails, child names, diagnoses, diary text, or other free-text fields. Candidate events are `registration_started` (role only), `registration_completed` (role only), `child_profile_created`, `specialist_assignment_created`, `first_session_created`, `first_diary_entry_created` (event only), `recovery_requested`, and a privacy-reviewed return-use event. Measure completion between steps and recurring use by role; do not record content to explain drop-off.

## Deliberately deferred

Do not add generative AI, in-app messaging, video, native apps, school integrations, or child-facing gamification before the core guardian–specialist workflow is observed with users. Notifications should be validated separately because lock-screen text can expose sensitive information.
