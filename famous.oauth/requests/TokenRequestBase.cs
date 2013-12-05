using famous.oauth.core.request;

namespace famous.oauth.requests
{
  /// <summary>
  /// OAuth 2.0 request for an access token as specified in http://tools.ietf.org/html/rfc6749#section-4.
  /// </summary>
  internal class TokenRequestBase
  {
    /// <summary>
    /// Gets or sets space-separated list of scopes as specified in http://tools.ietf.org/html/rfc6749#section-3.3.
    /// </summary>
    [HttpRequestParameter("scope")]
    public string Scope { get; set; }

    /// <summary>
    /// Gets or sets the Grant type. Sets <c>authorization_code</c> or <c>password</c> or <c>client_credentials</c> 
    /// or <c>refresh_token</c> or absolute URI of the extension grant type.
    /// </summary>
    [HttpRequestParameter("grant_type")]
    public string GrantType { get; set; }

    /// <summary>Gets or sets the client Identifier.</summary>
    [HttpRequestParameter("client_id")]
    public string ClientId { get; set; }

    /// <summary>Gets or sets the client Secret.</summary>
    [HttpRequestParameter("client_secret")]
    public string ClientSecret { get; set; }
  }
}