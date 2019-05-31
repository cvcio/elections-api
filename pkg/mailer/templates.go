package mailer

var (
	msgInvitation = `
<div style="font-family:'Segoe UI',Arial,Sans-Serif;font-size:10pt;">
    <p>
        Hi there,
    </p>
    <p>
        %s %s (%s) has invited you to join <a href="https://mediawatch.io">MediaWatch</a>.
    </p>
    <p>
        If you are interested you can create an account by following this URL:<br/>
        <a href="https://demo.mediawatch.io/auth/register">https://demo.mediawatch.io/auth/register</a>
    </p>
    <p>
        If you are are not interested just ignore this message.<br/>
        However, if you have any security or privacy discrimination concerns please contact us immediately by email: <a href="mailto:press@mediawatch.io">press@mediawatch.io</a>.
    </p>
    <p>
        Best regards,<br />
        MediaWatch team
    </p>
</div>
`
	msgNewPass = `
<div style="font-family:'Segoe UI',Arial,Sans-Serif;font-size:10pt;">
    <p>
        Hi %s,
    </p>
    <p>
        Your new password is: %s.
    </p>
    <p>
        To login to your account follow this URL:<br/>
        <a href="https://demo.mediawatch.io/auth/login">https://demo.mediawatch.io/auth/login</a>
    </p>
    <p>
        Best regards,<br />
        MediaWatch team
    </p>
</div>
`

	msgReset = `
<div style="font-family:'Segoe UI',Arial,Sans-Serif;font-size:10pt;">
    <p>
        Hi %s,
    </p>
    <p>
        You have requested a password reset for your MediaWatch account.<br/>
        Your 4-Digit Verification Code is: %d.
    </p>
    <p>
        To complete password reset follow this URL:<br/>
        <a href="https://demo.mediawatch.io/auth/reset/verify/%s">https://demo.mediawatch.io/auth/reset/verify/%s</a>
    </p>
    <p>
        The password reset link will expire in 24 hours.
    </p>
    <p>
        If you didn't request a password reset or made it accidentially just ignore this message.<br/>
        However, if you have any security concerns please contact us immediately by email: <a href="mailto:press@mediawatch.io">press@mediawatch.io</a>.
    </p>
    <p>
        Best regards,<br />
        MediaWatch team
    </p>
</div>
`

	msgPin = `
<div style="font-family:'Segoe UI',Arial,Sans-Serif;font-size:10pt;">
    <p>
        Hi %s,
    </p>
    <p>
        Your 4-Digit Verification Code is: %d.
    </p>
    <p>
        To complete your registration follow this URL:<br/>
        <a href="https://demo.mediawatch.io/auth/verify/%s">https://demo.mediawatch.io/auth/verify/%s</a>
    </p>
    <p>
        Please confirm your account within the next 24 hours.
    </p>
    <p>
        If you didn't register an account or made it accidentially just ignore this message.<br/>
        However, if you have any security concerns please contact us immediately by email: <a href="mailto:press@mediawatch.io">press@mediawatch.io</a>.
    </p>
    <p>
        Best regards,<br />
        MediaWatch team
    </p>
</div>
`

	msgDefault = `
%s<br/><br/>

Best regards,<br/>
MediaWatch team<br/>
`
)
