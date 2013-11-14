using System;
using System.Collections.Generic;
using System.Linq;
using famous.oauth.core;

namespace famous.oauth
{
  public class OAuth2Context
  {

    /// <summary>Gets the token server URL.</summary>
    public string TokenServerUrl { get; private set; }

    /// <summary>Gets or sets the authorization server URL.</summary>
    public string AuthorizationServerUrl { get; private set; }

    /// <summary>Gets or sets the client secrets which includes the client identifier and its secret.</summary>
    public ClientSecrets ClientSecrets { get; set; }

    /// <summary>Gets or sets the data store used to store the token response.</summary>
    public IDataStorage DataStore { get; set; }

    /// <summary>
    /// Gets or sets the scopes which indicate the API access your application is requesting.
    /// </summary>
    public IEnumerable<string> Scopes { get; set; }

    /// <summary>
    /// Gets or sets the clock. The clock is used to determine if the token has expired, if so we will try to
    /// refresh it. The default value is <seealso cref="DateTime.Now"/>.
    /// </summary>
    public DateTime Clock { get; set; }

    /// <summary>Constructs a new initializer.</summary>
    /// <param name="authorizationServerUrl">Authorization server URL</param>
    /// <param name="tokenServerUrl">Token server URL</param>
    /// <param name="redirectUrl">Registered redirect URL</param>
    public OAuth2Context(string authorizationServerUrl, string tokenServerUrl, string redirectUrl)
    {
      AuthorizationServerUrl = authorizationServerUrl;
      TokenServerUrl = tokenServerUrl;
      Clock = DateTime.Now;
      RedirectUrl = redirectUrl;
      Scopes = Enumerable.Empty<string>();
    }

    public string RedirectUrl { get; set; }
  }
}