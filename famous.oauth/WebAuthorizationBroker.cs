using System;
using System.Threading;
using System.Threading.Tasks;
using famous.oauth.core;
using famous.oauth.responses;

namespace famous.oauth
{
  /// <summary>A helper utility to manage the authorization code flow.</summary>
  public class WebAuthorizationBroker
  {
    private readonly AuthorizationCodeFlow flow;

    private readonly ICodeReceiver code_receiver = new LocalServerCodeReceiver();
    public async Task<TokenResponse> AuthorizeAsync(string userid_hint, bool force, CancellationToken canceltoken)
    {
      if (string.IsNullOrEmpty(userid_hint))
      {
        throw new ArgumentException("can not be empty", userid_hint);
      }
      var token = force ? null : await flow.LoadTokenAsync(userid_hint, canceltoken).ConfigureAwait(false) ;

      if (token != null && (token.RefreshToken != null || !token.IsExpired(DateTime.Now))) return token;
      var code_req = flow.MakeAuthorizationCodeRequest(code_receiver.CallbackUrl);
      var code_resp = await code_receiver.ReceiveCodeAsync(code_req.ToString(), canceltoken);
     
      if (string.IsNullOrEmpty(code_resp.Code))
      {
        throw new ResponseException<TokenErrorResponse>(new TokenErrorResponse()
        { 
          Error = code_resp.Error,
          ErrorDescription = code_resp.ErrorDescription,
          ErrorUri = code_resp.ErrorUri,
        });
      }
      token =
        await flow.ExchangeCodeForTokenAsync(userid_hint, code_resp.Code, canceltoken)
          .ConfigureAwait(false);
      return token;
    }

    public WebAuthorizationBroker(OAuth2Context c)
    {
      flow = new AuthorizationCodeFlow(c);
    }
  }
}