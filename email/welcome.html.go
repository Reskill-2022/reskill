package email

const welcomeHTML = `
<!DOCTYPE html PUBLIC "-//W3C//DTD XHTML 1.0 Transitional//EN" "http://www.w3.org/TR/xhtml1/DTD/xhtml1-transitional.dtd">
<html xmlns="http://www.w3.org/1999/xhtml" xmlns:v="urn:schemas-microsoft-com:vml" xmlns:o="urn:schemas-microsoft-com:office:office">
<head>
  <meta http-equiv="Content-Type" content="text/html; charset=utf-8">
  <!--[if !mso]><!-->
  <meta http-equiv="X-UA-Compatible" content="IE=edge">
  <!--<![endif]-->
  <meta name="viewport" content="width=device-width, initial-scale=1.0">
  <meta name="format-detection" content="telephone=no">
  <meta name="x-apple-disable-message-reformatting">
  <title></title>
  <style type="text/css">
    @media screen {
      @font-face {
        font-family: 'Lato';
        font-style: normal;
        font-weight: 300;
        src: url('https://fonts.gstatic.com/s/lato/v23/S6u9w4BMUTPHh7USSwaPHw.woff') format('woff'), url('https://fonts.gstatic.com/s/lato/v23/S6u9w4BMUTPHh7USSwaPGQ.woff2') format('woff2');
      }
      @font-face {
        font-family: 'Lato';
        font-style: normal;
        font-weight: 400;
        src: url('https://fonts.gstatic.com/s/lato/v23/S6uyw4BMUTPHjxAwWA.woff') format('woff'), url('https://fonts.gstatic.com/s/lato/v23/S6uyw4BMUTPHjxAwXg.woff2') format('woff2');
      }
      @font-face {
        font-family: 'Lato';
        font-style: normal;
        font-weight: 700;
        src: url('https://fonts.gstatic.com/s/lato/v23/S6u9w4BMUTPHh6UVSwaPHw.woff') format('woff'), url('https://fonts.gstatic.com/s/lato/v23/S6u9w4BMUTPHh6UVSwaPGQ.woff2') format('woff2');
      }
    }
  </style>
  <style type="text/css">
    #outlook a {
      padding: 0;
    }

    .ReadMsgBody,
    .ExternalClass {
      width: 100%;
    }

    .ExternalClass,
    .ExternalClass p,
    .ExternalClass td,
    .ExternalClass div,
    .ExternalClass span,
    .ExternalClass font {
      line-height: 100%;
    }

    div[style*="margin: 14px 0"],
    div[style*="margin: 16px 0"] {
      margin: 0 !important;
    }

    table,
    td {
      mso-table-lspace: 0;
      mso-table-rspace: 0;
    }

    table,
    tr,
    td {
      border-collapse: collapse;
    }

    body,
    td,
    th,
    p,
    div,
    li,
    a,
    span {
      -webkit-text-size-adjust: 100%;
      -ms-text-size-adjust: 100%;
      mso-line-height-rule: exactly;
    }

    img {
      border: 0;
      outline: none;
      line-height: 100%;
      text-decoration: none;
      -ms-interpolation-mode: bicubic;
    }

    a[x-apple-data-detectors] {
      color: inherit !important;
      text-decoration: none !important;
    }

    body {
      margin: 0;
      padding: 0;
      width: 100% !important;
      -webkit-font-smoothing: antialiased;
    }

    .pc-gmail-fix {
      display: none;
      display: none !important;
    }

    @media screen and (min-width: 621px) {
      .pc-email-container {
        width: 620px !important;
      }
    }
  </style>
  <style type="text/css">
    @media screen and (max-width:620px) {
      .pc-sm-mw-100pc {
        max-width: 100% !important
      }
      .pc-sm-p-25-10-15 {
        padding: 25px 10px 15px !important
      }
      .pc-sm-p-24-20-30 {
        padding: 24px 20px 30px !important
      }
    }
  </style>
  <style type="text/css">
    @media screen and (max-width:525px) {
      .pc-xs-w-100pc {
        width: 100% !important
      }
      .pc-xs-p-10-0-0 {
        padding: 10px 0 0 !important
      }
      .pc-xs-p-15-0-5 {
        padding: 15px 0 5px !important
      }
      .pc-xs-br-disabled br {
        display: none !important
      }
      .pc-xs-p-15-10-20 {
        padding: 15px 10px 20px !important
      }
      .pc-xs-h-100 {
        height: 100px !important
      }
      .pc-xs-fs-30 {
        font-size: 30px !important
      }
      .pc-xs-lh-42 {
        line-height: 42px !important
      }
    }
  </style>
  <!--[if mso]>
    <style type="text/css">
        .pc-fb-font {
            font-family: Helvetica, Arial, sans-serif !important;
        }
    </style>
    <![endif]-->
  <!--[if gte mso 9]><xml><o:OfficeDocumentSettings><o:AllowPNG/><o:PixelsPerInch>96</o:PixelsPerInch></o:OfficeDocumentSettings></xml><![endif]-->
</head>
<body style="width: 100% !important; margin: 0; padding: 0; mso-line-height-rule: exactly; -webkit-font-smoothing: antialiased; -webkit-text-size-adjust: 100%; -ms-text-size-adjust: 100%; background-color: #f4f4f4; background-position: center center; background-repeat: no-repeat; background-size: cover" class="" data-new-gr-c-s-check-loaded="14.1078.0" data-gr-ext-installed="">
  <div style="display: none !important; visibility: hidden; opacity: 0; overflow: hidden; mso-hide: all; height: 0; width: 0; max-height: 0; max-width: 0; font-size: 1px; line-height: 1px; color: #151515;">This is preheader text. Some</div>
  <div style="display: none !important; visibility: hidden; opacity: 0; overflow: hidden; mso-hide: all; height: 0; width: 0; max-height: 0; max-width: 0; font-size: 1px; line-height: 1px;">
    ‌&nbsp;‌&nbsp;‌&nbsp;‌&nbsp;‌&nbsp;‌&nbsp;‌&nbsp;‌&nbsp;‌&nbsp;‌&nbsp;‌&nbsp;‌&nbsp;‌&nbsp;‌&nbsp;‌&nbsp;‌&nbsp;‌&nbsp;‌&nbsp;‌&nbsp;‌&nbsp;‌&nbsp;‌&nbsp;‌&nbsp;‌&nbsp;‌&nbsp;‌&nbsp;‌&nbsp;‌&nbsp;‌&nbsp;‌&nbsp;‌&nbsp;‌&nbsp;‌&nbsp;‌&nbsp;‌&nbsp;‌&nbsp;‌&nbsp;‌&nbsp;‌&nbsp;‌&nbsp;‌&nbsp;‌&nbsp;‌&nbsp;‌&nbsp;‌&nbsp;‌&nbsp;‌&nbsp;‌&nbsp;‌&nbsp;‌&nbsp;‌&nbsp;‌&nbsp;‌&nbsp;‌&nbsp;‌&nbsp;‌&nbsp;‌&nbsp;‌&nbsp;‌&nbsp;‌&nbsp;‌&nbsp;‌&nbsp;‌&nbsp;‌&nbsp;‌&nbsp;‌&nbsp;‌&nbsp;‌&nbsp;‌&nbsp;‌&nbsp;‌&nbsp;‌&nbsp;‌&nbsp;‌&nbsp;‌&nbsp;‌&nbsp;‌&nbsp;‌&nbsp;‌&nbsp;‌&nbsp;‌&nbsp;‌&nbsp;‌&nbsp;‌&nbsp;‌&nbsp;‌&nbsp;‌&nbsp;‌&nbsp;‌&nbsp;‌&nbsp;‌&nbsp;‌&nbsp;‌&nbsp;‌&nbsp;‌&nbsp;‌&nbsp;‌&nbsp;‌&nbsp;‌&nbsp;‌&nbsp;‌&nbsp;‌&nbsp;‌&nbsp;‌&nbsp;‌&nbsp;‌&nbsp;‌&nbsp;‌&nbsp;‌&nbsp;‌&nbsp;‌&nbsp;‌&nbsp;‌&nbsp;‌&nbsp;‌&nbsp;‌&nbsp;‌&nbsp;‌&nbsp;‌&nbsp;‌&nbsp;‌&nbsp;‌&nbsp;‌&nbsp;‌&nbsp;‌&nbsp;‌&nbsp;‌&nbsp;‌&nbsp;‌&nbsp;‌&nbsp;‌&nbsp;‌&nbsp;‌&nbsp;‌&nbsp;‌&nbsp;‌&nbsp;‌&nbsp;‌&nbsp;‌&nbsp;‌&nbsp;‌&nbsp;‌&nbsp;‌&nbsp;‌&nbsp;‌&nbsp;‌&nbsp;‌&nbsp;‌&nbsp;‌&nbsp;‌&nbsp;‌&nbsp;‌&nbsp;‌&nbsp;‌&nbsp;‌&nbsp;‌&nbsp;
  </div>
  <table class="pc-email-body" width="100%" bgcolor="#f4f4f4" border="0" cellpadding="0" cellspacing="0" role="presentation" style="table-layout: fixed;">
    <tbody>
      <tr>
        <td class="pc-email-body-inner" align="center" valign="top">
          <!--[if gte mso 9]>
            <v:background xmlns:v="urn:schemas-microsoft-com:vml" fill="t">
                <v:fill type="tile" src="" color="#f4f4f4"></v:fill>
            </v:background>
            <![endif]-->
          <!--[if (gte mso 9)|(IE)]><table width="620" align="center" border="0" cellspacing="0" cellpadding="0" role="presentation"><tr><td width="620" align="center" valign="top"><![endif]-->
          <table class="pc-email-container" width="100%" align="center" border="0" cellpadding="0" cellspacing="0" role="presentation" style="margin: 0 auto; max-width: 620px;">
            <tbody>
              <tr>
                <td align="left" valign="top" style="padding: 0 10px;">
                  <table width="100%" border="0" cellpadding="0" cellspacing="0" role="presentation">
                    <tbody>
                      <tr>
                        <td height="20" style="font-size: 1px; line-height: 1px;">&nbsp;</td>
                      </tr>
                    </tbody>
                  </table>
    
                  <!-- BEGIN MODULE: Content 2 -->
                  <table border="0" cellpadding="0" cellspacing="0" width="100%" role="presentation">
                    <tbody>
                      <tr>
                        <td class="pc-sm-p-25-10-15 pc-xs-p-15-0-5" valign="top" bgcolor="#ffffff" style="padding: 30px 20px 20px; background-color: #ffffff; border-radius: 8px">
                          <table border="0" cellpadding="0" cellspacing="0" width="100%" role="presentation">
                            <tbody>
                              <tr>
                                <img style="    width: 51%;
                                align-items: center;
                                display: flex;
                               
                                margin: 15px auto;
                                " src="https://res.cloudinary.com/dgc6ootad/image/upload/v1662062691/Reskill%20form/Logo_2_g1trhg_1_yraemb.svg" alt="">
                                <td valign="top" style="padding: 0 20px;">
                                  <table border="0" cellpadding="0" cellspacing="0" width="100%" role="presentation">
                                    <tbody>
                                      <tr>
                                        <td class="pc-fb-font" valign="top" style="padding: 10px 0 0; font-family: 'Lato', Helvetica, Arial, sans-serif; font-size: 23px; text-align:center; font-weight: 200;padding-bottom: 15px; line-height: 34px; letter-spacing: -0.4px; color: #1d1e1f">INCREASING ACCESS TO CAREERS IN TECH.</td>
                                      </tr>
                                    </tbody>
                                  </table>
                                </td>
                              </tr>
                            </tbody>
                           
                            <tbody>
                              <tr>
                             
                                <td valign="top" style="padding: 0 20px;">
                                  <table border="0" cellpadding="0" cellspacing="0" width="100%" role="presentation">
                                    <tbody>
                                      <tr>
                                        <hr>
                                        <td class="pc-fb-font" valign="top" style="padding: 10px 0 0; font-family: 'Lato', Helvetica, Arial, sans-serif; font-size: 21px; font-weight: 700; line-height: 34px; letter-spacing: -0.4px; color: #2a6cdf">Welcome to Reskill Americans!</td>
                                      </tr>
                                    </tbody>
                                  </table>
                                </td>
                              </tr>
                            </tbody>
                            <tbody>
                              <tr>
                                <td class="pc-fb-font" valign="top" style="font-family: 'Lato', Helvetica, Arial, sans-serif; padding: 10px 20px 0; line-height: 28px; font-size: 18px; font-weight: 300; letter-spacing: -0.2px; color: #483b3b">
                                  <p>You are registered as {{ .Name }} at {{ .Email }}.<br><br>Thank you for completing our enrollment form with your information. On October 3, 2022 you will begin your exciting journey towards a career in software development! We look forward to introducing you to new skills that will help prepare you for a starting position in the tech world.<br><br><strong>Next Steps:</strong> Prior to October 3rd, you will receive login details for our program’s Learning Management System (LMS) and our online chatrooms. There is nothing you need to do before then to prepare (except to respond to any inquiries we might send!) In the meantime, please visit our FAQ page on our website<br>and read the key information below about our program to familiarize yourself.<br><br><strong>WHAT YOU NEED TO KNOW:<br><br></strong>What You Need: In addition to a passion to learn and determination to keep studying with us, you’ll need a laptop, tablet, or desktop PC with an internet connection. It is possible to do every part of this program using all online resources, which is why you only need a computer with a web browser. All the software used will be free, either online or installed on your computer. For more details about tech specs, see our FAQ page.<br><br><strong>Time Commitment:</strong> This is a seven-month, 100% online/remote learning experience. Our program is designed to fit into your schedule so that you can choose your own hours, pending the schedules of any peers with whom you might be working. You will need to dedicate at least 15 hours per week to be successful.<br><br><strong>How It Works:</strong> You will learn the fundamental concepts of software development and product design through a range of formats: recorded and live video sessions, training modules, and online chats or video conferences with your instructors, mentors, and peers.<br><br><strong>We Are Here to Help:</strong> Whatever your level of knowledge is about the industry, our instructors and mentors will meet you at that level and work with you to address your questions. While it is not possible to become an expert in just seven months (expertise takes years of experience), we will give you everything you need to get started and the resources required to grow.<br><br>If you have additional questions, please contact us by hitting reply on this email.<br><br>Welcome aboard!<br>Regards,<br>The Reskill Americans Team</p>
                                </td>
                              </tr>
                              <tr>
                                <td height="4" style="font-size: 1px; line-height: 1px">&nbsp;</td>
                              </tr>
                            </tbody>
                            <tbody>
                            </tbody>
                          </table>
                        </td>
                      </tr>
                    </tbody>
                  </table>
                  <!-- END MODULE: Content 2 -->
                  <table width="100%" border="0" cellpadding="0" cellspacing="0" role="presentation">
                    <tbody>
                      <tr>
                        <td height="20" style="font-size: 1px; line-height: 1px;">&nbsp;</td>
                      </tr>
                    </tbody>
                  </table>
                </td>
              </tr>
            </tbody>
          </table>
          <!--[if (gte mso 9)|(IE)]></td></tr></table><![endif]-->
        </td>
      </tr>
    </tbody>
  </table>
  <!-- Fix for Gmail on iOS -->
  <div class="pc-gmail-fix" style="white-space: nowrap; font: 15px courier; line-height: 0;">&nbsp; &nbsp; &nbsp; &nbsp; &nbsp; &nbsp; &nbsp; &nbsp; &nbsp; &nbsp; &nbsp; &nbsp; &nbsp; &nbsp; &nbsp; &nbsp; &nbsp; &nbsp; &nbsp; &nbsp; &nbsp; &nbsp; &nbsp; &nbsp; &nbsp; &nbsp; &nbsp; &nbsp; &nbsp; &nbsp; </div>
</body>
<grammarly-desktop-integration data-grammarly-shadow-root="true"></grammarly-desktop-integration>
</html>
`
