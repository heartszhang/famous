using famous.oauth.core;
using famous.oauth.core.request;

namespace famous.oauth.requests
{
  /// <summary>
  /// OAuth 2.0 request to refresh an access token using a refresh token as specified in 
  /// http://tools.ietf.org/html/rfc6749#section-6.
  /// </summary>
  internal class RefreshTokenRequest : TokenRequestBase
  {
    /// <summary>Gets or sets the Refresh token issued to the client.</summary>
    [HttpRequestParameter("refresh_token")]
    public string RefreshToken { get; set; }

    /// <summary>
    /// Constructs a new refresh code token request and sets grant_type to <c>refresh_token</c>.
    /// </summary>
    public RefreshTokenRequest()
    {
      GrantType = "refresh_token";
    }
  }
}