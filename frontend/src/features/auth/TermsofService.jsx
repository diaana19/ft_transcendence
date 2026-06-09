// TermsOfService is a static page that renders the terms of service sections.
export default function TermsOfService() {
    return (
        <div className="w-full mx-auto px-6 py-8 bg-transparent">
            <div className="mb-8">
                <h1 className="text-2xl font-semibold mb-1" style={{ color: '#2c2c2a' }}>
                    Terms of Service
                </h1>
                <p className="text-sm" style={{ color: '#b4b2a9' }}>
                    Effective date: June 1, 2026 · Last revised: June 1, 2026
                </p>
            </div>

            {[
                {
                    title: '1. Acceptance of Terms',
                    text: `Welcome to Synk. These Terms of Service ("Terms") constitute a legally binding agreement between you ("User", "you", or "your") and Synk ("we", "us", or "our"), governing your access to and use of the Synk platform, including its website, mobile application, APIs, and all associated services (collectively, the "Service").

By creating an account, accessing, or using any part of the Service, you acknowledge that you have read, understood, and agree to be bound by these Terms in their entirety, as well as our Privacy Policy, which is incorporated herein by reference. If you do not agree to these Terms, you must immediately cease all use of the Service and delete your account.

These Terms apply to all visitors, registered users, and any other individuals who access or use the Service in any manner. We reserve the right to modify these Terms at any time at our sole discretion. Your continued use of the Service following any modification constitutes your acceptance of the revised Terms. It is your responsibility to review these Terms periodically.

Synk is an academic project developed as part of the ft_transcendence curriculum at École 42. While it is developed in an academic context, all terms herein are intended to be taken seriously and reflect best practices for social platform governance.`,
                },
                {
                    title: '2. Eligibility and Account Registration',
                    text: `To register for and use the Service, you must be at least 16 years of age. By creating an account, you represent and warrant that: (a) you are at least 16 years old; (b) you have the legal capacity to enter into a binding agreement; (c) you are not prohibited from using the Service under applicable law; and (d) all information you provide during registration is accurate, current, and complete.

You agree to maintain the accuracy of your registration information and to promptly update it if it changes. You are solely responsible for maintaining the confidentiality of your account credentials, including your password. You must notify us immediately upon becoming aware of any unauthorized access to or use of your account.

You may not create an account on behalf of another person without their express written consent. You may not use a username that impersonates another person or entity, is offensive, or violates any third-party rights. We reserve the right to reject, suspend, or terminate any account at our sole discretion.

Each user is permitted to maintain only one active account. Creating multiple accounts for the purpose of circumventing restrictions, bans, or limitations is strictly prohibited and may result in permanent termination of all associated accounts.`,
                },
                {
                    title: '3. User Content and License',
                    text: `The Service allows you to create, upload, submit, store, and share content including but not limited to text posts, comments, images, videos, profile information, and direct messages ("User Content"). You retain full ownership of all User Content you submit to the Service.

By submitting User Content, you grant Synk a worldwide, non-exclusive, royalty-free, sublicensable, and transferable license to use, reproduce, distribute, prepare derivative works of, display, and perform your User Content in connection with the operation and provision of the Service. This license is solely for the purpose of operating the platform and will terminate when you delete your content or account, subject to reasonable technical delays for content removal from backup systems.

You represent and warrant that: (a) you own or have the necessary rights to the User Content you submit; (b) your User Content does not infringe the intellectual property, privacy, publicity, or other rights of any third party; (c) your User Content complies with these Terms and all applicable laws and regulations.

You acknowledge that Synk does not pre-screen User Content but reserves the right (but not the obligation) to review, remove, or modify any User Content at any time and for any reason, including content that we determine, in our sole discretion, violates these Terms or is otherwise objectionable.`,
                },
                {
                    title: '4. Prohibited Conduct',
                    text: `You agree not to engage in any of the following prohibited activities in connection with your use of the Service:

(a) Posting, transmitting, or distributing content that is unlawful, defamatory, obscene, pornographic, invasive of another's privacy, hateful, or racially, ethnically, or otherwise objectionable;
(b) Harassing, threatening, intimidating, stalking, or otherwise abusing or harming another user or any third party;
(c) Impersonating any person or entity or misrepresenting your affiliation with any person or entity;
(d) Engaging in any form of deceptive, fraudulent, or misleading conduct;
(e) Attempting to gain unauthorized access to any portion of the Service, other user accounts, or any systems or networks connected to the Service;
(f) Using automated tools, bots, scrapers, or similar mechanisms to access, collect, or extract data from the Service without our express written permission;
(g) Interfering with or disrupting the integrity or performance of the Service or the servers or networks connected to the Service;
(h) Uploading or transmitting viruses, malware, or any other malicious code;
(i) Using the Service for any commercial purpose without our prior written consent;
(j) Circumventing, disabling, or otherwise interfering with security-related features of the Service;
(k) Violating any applicable local, national, or international laws or regulations.

Violation of any of the above may result in immediate suspension or termination of your account and, where applicable, referral to law enforcement authorities.`,
                },
                {
                    title: '5. Intellectual Property',
                    text: `The Service and its original content (excluding User Content), features, functionality, design, logos, trademarks, and software are and will remain the exclusive property of Synk and its licensors. The Service is protected by applicable intellectual property laws in France and internationally.

Our trademarks and trade dress may not be used in connection with any product or service without the prior written consent of Synk. Nothing in these Terms grants you any right to use the Synk name, logo, or other identifying marks.

You acknowledge that we may use feedback, suggestions, or ideas you provide about the Service without any obligation to compensate you, and that such feedback may be incorporated into the Service at our sole discretion.`,
                },
                {
                    title: '6. Privacy and Data Protection',
                    text: `Your use of the Service is subject to our Privacy Policy, which is incorporated into these Terms by reference. By using the Service, you consent to the collection, use, and processing of your personal data as described in our Privacy Policy, in accordance with the General Data Protection Regulation (GDPR) and applicable French data protection law.

We are committed to protecting the privacy of users under the age of 18. If you become aware that a minor under the age of 16 has provided us with personal data, please contact us immediately and we will take steps to remove that information.`,
                },
                {
                    title: '7. Third-Party Services',
                    text: `The Service may integrate with or allow access through third-party authentication services, including but not limited to GitHub OAuth. Your use of such third-party services is subject to the respective terms and privacy policies of those third parties. We have no control over, and assume no responsibility for, the content, privacy policies, or practices of any third-party services.

The Service may contain links to third-party websites or services. These links are provided for your convenience only. We have no control over the content of those sites or resources, and we accept no responsibility for them or for any loss or damage that may arise from your use of them.`,
                },
                {
                    title: '8. Termination',
                    text: `You may terminate your account at any time by using the account deletion feature available in your profile settings. Upon termination, your right to use the Service will immediately cease.

We reserve the right to suspend or terminate your account and access to the Service at any time and for any reason, including but not limited to: violation of these Terms, conduct that we believe is harmful to other users or third parties, extended periods of inactivity, or if we are required to do so by law. We will make reasonable efforts to notify you prior to termination, except where prohibited by law or where immediate action is necessary to protect the Service or its users.

Upon termination of your account for any reason, we will process the deletion of your personal data in accordance with our Privacy Policy and applicable law. Provisions of these Terms that by their nature should survive termination shall survive, including but not limited to intellectual property provisions, disclaimers, and limitations of liability.`,
                },
                {
                    title: '9. Disclaimers and Limitation of Liability',
                    text: `THE SERVICE IS PROVIDED ON AN "AS IS" AND "AS AVAILABLE" BASIS WITHOUT WARRANTIES OF ANY KIND, EITHER EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE, NON-INFRINGEMENT, OR COURSE OF PERFORMANCE.

SYNK AND ITS DEVELOPERS DO NOT WARRANT THAT: (A) THE SERVICE WILL FUNCTION UNINTERRUPTED, SECURE, OR AVAILABLE AT ANY PARTICULAR TIME OR LOCATION; (B) ANY ERRORS OR DEFECTS WILL BE CORRECTED; (C) THE SERVICE IS FREE OF VIRUSES OR OTHER HARMFUL COMPONENTS; OR (D) THE RESULTS OF USING THE SERVICE WILL MEET YOUR REQUIREMENTS.

TO THE FULLEST EXTENT PERMITTED BY APPLICABLE LAW, SYNK SHALL NOT BE LIABLE FOR ANY INDIRECT, INCIDENTAL, SPECIAL, CONSEQUENTIAL, OR PUNITIVE DAMAGES, INCLUDING WITHOUT LIMITATION, LOSS OF PROFITS, DATA, USE, GOODWILL, OR OTHER INTANGIBLE LOSSES, RESULTING FROM YOUR ACCESS TO OR USE OF (OR INABILITY TO ACCESS OR USE) THE SERVICE.`,
                },
                {
                    title: '10. Governing Law and Dispute Resolution',
                    text: `These Terms shall be governed by and construed in accordance with the laws of France, without regard to its conflict of law provisions. You agree that any dispute, controversy, or claim arising out of or relating to these Terms, or the breach, termination, or invalidity thereof, shall be subject to the exclusive jurisdiction of the competent courts of France.

If any provision of these Terms is found to be unenforceable or invalid under applicable law, such provision shall be modified to the minimum extent necessary to make it enforceable, and the remaining provisions shall continue in full force and effect.

These Terms, together with our Privacy Policy and any other legal notices published by us on the Service, constitute the entire agreement between you and Synk with respect to the Service and supersede all prior agreements, understandings, and negotiations.`,
                },
                {
                    title: '11. Contact Information',
                    text: `If you have any questions, concerns, or requests regarding these Terms of Service, please contact us through the Synk application or via the École 42 platform. We will make reasonable efforts to respond to all inquiries in a timely manner.

For requests related to your personal data, including data access, correction, or deletion, please use the GDPR tools available in your account settings or contact us directly as described in our Privacy Policy.`,
                },
            ].map((section) => (
                <div
                    key={section.title}
                    className="mb-6 rounded-xl px-4 py-4"
                    style={{ background: 'white', border: '0.5px solid #ede8fd' }}
                >
                    <div className="text-sm font-semibold mb-3" style={{ color: '#534ab7' }}>
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
