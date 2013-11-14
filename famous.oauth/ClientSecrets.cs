namespace famous.oauth
{
  /// <summary>Client credential details for installed and web applications.</summary>
  public sealed class ClientSecrets
  {
    /// <summary>Gets or sets the client identifier.</summary>
    [Newtonsoft.Json.JsonProperty("client_id")]
    public string ClientId { get; set; }

    /// <summary>Gets or sets the client Secret.</summary>
    [Newtonsoft.Json.JsonProperty("client_secret")]
    public string ClientSecret { get; set; }
  }
}