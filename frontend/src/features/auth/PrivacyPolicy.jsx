export default function PrivacyPolicy() {
    return (
        <div className="w-full mx-auto px-6 py-8 bg-transparent">
            <div className="mb-8">
                <h1 className="text-2xl font-semibold mb-1" style={{ color: '#2c2c2a' }}>
                    Privacy Policy
                </h1>
                <p className="text-sm" style={{ color: '#b4b2a9' }}>
                    Effective date: June 1, 2026 · Last revised: June 1, 2026
                </p>
            </div>

            {[
                {
                    title: '1. Introduction and Data Controller',
                    text: `This Privacy Policy describes how Synk ("we", "us", or "our") collects, uses, stores, and protects your personal data when you use the Synk platform and its associated services (the "Service"). This policy is written in accordance with the General Data Protection Regulation (GDPR — Regulation EU 2016/679) and applicable French data protection legislation.

Synk is an academic social platform developed as part of the ft_transcendence curriculum at École 42. For the purposes of data protection law, the data controller is the Synk project team operating under the academic framework of École 42, France.

By using the Service, you acknowledge that you have read and understood this Privacy Policy. If you do not agree with our data practices described herein, you should not use the Service. We encourage you to read this policy carefully and in its entirety before using the platform.`,
                },
                {
                    title: '2. Data We Collect',
                    text: `We collect the following categories of personal data:

Account and Registration Data: When you register for Synk, we collect your username, email address, and a hashed version of your password. If you register via GitHub OAuth, we may receive your GitHub username, email address, and public profile information as provided by GitHub.

Profile Data: Information you voluntarily provide to personalize your account, including your display name, biography, profile picture (avatar), and banner image (wallpaper).

Content Data: All content you create or submit through the Service, including posts, comments, direct messages, media files (images and videos), and any associated metadata.

Interaction Data: Records of your interactions with other users and content on the platform, including likes, follows, friend connections, friend requests, and notifications.

Authentication and Security Data: Data related to your authentication sessions, including JWT tokens, session identifiers, two-factor authentication (2FA) configuration and status, and login timestamps.

Gamification Data: Computed statistics derived from your activity on the platform, including post counts, like counts, follower and following counts, activity scores, and earned badges.

Technical Data: When you access the Service, we may automatically collect technical information including your IP address, browser type and version, operating system, device identifiers, and access timestamps for security and operational purposes.`,
                },
                {
                    title: '3. Legal Basis for Processing',
                    text: `We process your personal data on the following legal bases under GDPR:

Performance of a Contract (Article 6(1)(b)): Processing necessary to provide you with the Service, including account management, content delivery, messaging functionality, and social features.

Legitimate Interests (Article 6(1)(f)): Processing necessary for our legitimate interests in operating, securing, and improving the Service, including fraud prevention, abuse detection, and platform analytics, provided such interests are not overridden by your fundamental rights and freedoms.

Consent (Article 6(1)(a)): Where we rely on your consent for specific processing activities, such as optional features or communications. You may withdraw your consent at any time without affecting the lawfulness of processing prior to withdrawal.

Legal Obligation (Article 6(1)(c)): Processing necessary to comply with applicable legal obligations, including data retention requirements and responses to lawful requests from public authorities.`,
                },
                {
                    title: '4. How We Use Your Data',
                    text: `We use your personal data for the following purposes:

To provide and operate the Service: Creating and managing your account, authenticating your identity, enabling social features (posts, comments, messaging, follows, friendships), delivering notifications, and computing gamification statistics.

To ensure security and prevent abuse: Detecting and preventing fraudulent activity, unauthorized access, spam, and violations of our Terms of Service. This includes monitoring for suspicious login activity and enforcing rate limits.

To improve the Service: Analyzing usage patterns and performance data to identify and fix bugs, improve features, and enhance the overall user experience.

To comply with legal obligations: Responding to lawful requests from public authorities, complying with data retention requirements, and fulfilling our obligations under applicable law.

We do not use your personal data for automated decision-making or profiling that produces legal or similarly significant effects. We do not sell, rent, or otherwise commercially exploit your personal data to third parties.`,
                },
                {
                    title: '5. Data Sharing and Third Parties',
                    text: `We do not sell or rent your personal data to any third party. We may share your data only in the following limited circumstances:

Service Providers: We may engage third-party service providers who process data on our behalf strictly for the purpose of operating the Service (e.g., hosting infrastructure). Such providers are bound by data processing agreements and may only process your data in accordance with our instructions.

GitHub OAuth: If you choose to authenticate via GitHub, your authentication is handled by GitHub in accordance with GitHub's Privacy Policy. We receive only the data necessary to create and manage your Synk account (username, email, and public profile data).

Legal Compliance: We may disclose your data to law enforcement, regulatory authorities, or other third parties if required to do so by applicable law, court order, or governmental regulation, or if we believe in good faith that such disclosure is necessary to protect our legal rights, prevent fraud or illegal activity, or ensure the safety of our users.

Business Transfers: In the event of a merger, acquisition, or asset transfer, your data may be transferred as part of that transaction, subject to the acquirer's agreement to comply with this Privacy Policy.`,
                },
                {
                    title: '6. Data Retention',
                    text: `We retain your personal data for as long as your account remains active or as necessary to provide you with the Service. Specifically:

Account and profile data is retained for the duration of your account's existence on the platform. When you delete your account, your personal data is permanently deleted from our active systems within a reasonable timeframe, subject to the exceptions noted below.

User Content (posts, comments, messages, uploaded media) is deleted upon account deletion, unless it has been shared with other users in a way that makes independent deletion technically complex, in which case it will be anonymized.

Technical and security logs may be retained for a limited period (generally no longer than 90 days) for security, fraud prevention, and debugging purposes, even after account deletion.

Certain data may be retained for longer periods where required by applicable law, to resolve disputes, enforce our agreements, or comply with legal obligations. In such cases, we will limit processing of such data to the extent necessary for the applicable legal purpose.`,
                },
                {
                    title: '7. Your Rights Under GDPR',
                    text: `As a data subject under the GDPR, you have the following rights with respect to your personal data:

Right of Access (Article 15): You have the right to obtain confirmation of whether we process your personal data and, if so, to receive a copy of that data along with information about how it is processed. You can exercise this right using the data export feature available in your account settings.

Right to Rectification (Article 16): You have the right to request correction of inaccurate or incomplete personal data we hold about you. You can update most of your information directly in your account settings.

Right to Erasure ("Right to be Forgotten") (Article 17): You have the right to request deletion of your personal data in certain circumstances, including where the data is no longer necessary for the purposes for which it was collected. You can exercise this right by deleting your account through the platform settings.

Right to Restriction of Processing (Article 18): You have the right to request that we restrict processing of your personal data in certain circumstances.

Right to Data Portability (Article 20): You have the right to receive your personal data in a structured, commonly used, and machine-readable format and to transmit that data to another controller. Use the GDPR data export feature in your account settings to download your data.

Right to Object (Article 21): You have the right to object to processing of your personal data where we rely on legitimate interests as the legal basis for processing.

Right to Withdraw Consent: Where processing is based on your consent, you have the right to withdraw that consent at any time without affecting the lawfulness of processing prior to withdrawal.

To exercise any of the above rights, please use the GDPR tools available in your account settings or contact us directly. We will respond to your request within 30 days as required by applicable law.`,
                },
                {
                    title: '8. Security Measures',
                    text: `We implement appropriate technical and organizational measures to protect your personal data against unauthorized access, accidental loss, destruction, or alteration. These measures include:

Passwords are stored using industry-standard cryptographic hashing algorithms and are never stored in plaintext. All communications between your browser and our servers are encrypted using HTTPS/TLS. Authentication sessions are managed using short-lived JWT tokens and HttpOnly cookies, mitigating common web security vulnerabilities. Two-factor authentication (2FA) is available to all users and strongly recommended. We implement rate limiting and other anti-abuse measures to protect against unauthorized access attempts.

However, no security system is impenetrable and we cannot guarantee the absolute security of your data. In the event of a personal data breach that is likely to result in a risk to your rights and freedoms, we will notify you and the relevant supervisory authority as required by GDPR.`,
                },
                {
                    title: '9. Cookies and Similar Technologies',
                    text: `Synk uses a limited number of cookies and similar technologies strictly necessary for the operation of the Service:

Authentication Cookie (auth_token): An HttpOnly, Secure session cookie used to manage your authenticated session when logging in via GitHub OAuth. This cookie is strictly necessary and cannot be disabled without affecting your ability to use the Service.

We do not use third-party advertising cookies, tracking pixels, or cross-site analytics tools. We do not use your data to serve targeted advertisements.

By using the Service, you consent to the use of strictly necessary cookies as described above. You may configure your browser to refuse cookies, but please note that this may affect the functionality of the Service.`,
                },
                {
                    title: "10. Children's Privacy",
                    text: `The Service is not directed to individuals under the age of 16. We do not knowingly collect personal data from children under 16. If you are a parent or guardian and believe that your child under the age of 16 has provided us with personal data without your consent, please contact us immediately.

Upon receiving such a request, we will take immediate steps to verify the information and, if confirmed, delete the relevant personal data from our systems as quickly as technically possible. We reserve the right to terminate any account that we discover to be held by a user under the age of 16.`,
                },
                {
                    title: '11. International Data Transfers',
                    text: `The Service is operated and hosted within the European Union. Your personal data is processed and stored on servers located within the EU and is subject to the protections of French law and the GDPR.

If, in the future, we engage service providers located outside the European Economic Area (EEA), we will ensure that appropriate safeguards are in place for any such transfer of personal data, including Standard Contractual Clauses or other mechanisms approved by the European Commission.`,
                },
                {
                    title: '12. Changes to This Privacy Policy',
                    text: `We may update this Privacy Policy from time to time to reflect changes in our data practices, legal requirements, or the features of the Service. When we make material changes, we will notify you by posting the updated policy on the Service and updating the "Last revised" date at the top of this page. In some cases, we may provide additional notice such as an in-app notification.

Your continued use of the Service after the effective date of the revised Privacy Policy constitutes your acceptance of the updated policy. We encourage you to review this policy periodically. If you do not agree to the revised policy, you must cease using the Service and delete your account.`,
                },
                {
                    title: '13. Contact and Supervisory Authority',
                    text: `If you have any questions, concerns, or requests regarding this Privacy Policy or our data practices, please contact us through the Synk application or via the École 42 platform. We will do our best to respond to all inquiries within 30 days.

If you believe that our processing of your personal data infringes your rights under the GDPR, you have the right to lodge a complaint with the competent supervisory authority. In France, the relevant authority is:

Commission Nationale de l'Informatique et des Libertés (CNIL)
3 Place de Fontenoy — TSA 80715
75334 Paris Cedex 07
France
www.cnil.fr`,
                },
            ].map((section) => (
                <div
                    key={section.title}
                    className="mb-6 rounded-xl px-4 py-4"
                    style={{ background: 'white', border: '0.5px solid #e8f0fd' }}
                >
                    <div className="text-sm font-semibold mb-3" style={{ color: '#185fa5' }}>
                        {section.title}
                    </div>
                    <div
                        className="text-xs leading-relaxed whitespace-pre-line"
                        style={{ color: '#5f5e5a' }}
                    >
                        {section.text}
                    </div>
                </div>
            ))}
        </div>
    )
}
