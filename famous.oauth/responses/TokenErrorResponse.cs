namespace famous.oauth.responses
{
  /// <summary>
  /// OAuth 2.0 model for a unsuccessful access token response as specified in 
  /// http://tools.ietf.org/html/rfc6749#section-5.2.
  /// </summary>
  public class TokenErrorResponse
  {
    /// <summary>
    /// Gets or sets error code (e.g. "invalid_request", "invalid_client", "invalid_grant", "unauthorized_client", 
    /// "unsupported_grant_type", "invalid_scope") as specified in http://tools.ietf.org/html/rfc6749#section-5.2.
    /// </summary>
    [Newtonsoft.Json.JsonProperty("error")]
    public string Error { get; set; }

    /// <summary>
    /// Gets or sets a human-readable text which provides additional information used to assist the client 
    /// developer in understanding the error occurred.
    /// </summary>
    [Newtonsoft.Json.JsonProperty("error_description")]
    public string ErrorDescription { get; set; }

    /// <summary>
    /// Gets or sets the URI identifying a human-readable web page with provides information about the error.
    /// </summary>
    [Newtonsoft.Json.JsonProperty("error_uri")]
    public string ErrorUri { get; set; }

    public override string ToString()
    {
      return string.Format("Error:\"{0}\", Description:\"{1}\", Uri:\"{2}\"", Error, ErrorDescription, ErrorUri);
    }
  }
}