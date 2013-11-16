using System;
using System.Collections.Generic;
using System.Net.Http;
using System.Threading;
using System.Threading.Tasks;
using famous.oauth.core;
using famous.oauth.requests;
using famous.oauth.responses;
using famous.oauth.utils;
using Newtonsoft.Json;
namespace famous.oauth
{
  /// <summary>
  /// Thread-safe OAuth 2.0 authorization code flow that manages and persists end-user credentials.
  /// <para>
  /// This is designed to simplify the flow in which an end-user authorizes the application to access their protected
  /// data, and then the application has access to their data based on an access token and a refresh token to refresh 
  /// that access token when it expires.
  /// </para>
  /// </summary>
  internal class AuthorizationCodeFlow 
  {
    /// <summary>An OAuth2Context class for the authorization code flow. </summary>
    private readonly OAuth2Context _;
    /// <summary>Constructs a new flow using the initializer's properties.</summary>
    public AuthorizationCodeFlow(OAuth2Context initializer)
    {
      _ = initializer;

      if (_.ClientSecrets == null || _.DataStore == null)
      {
        throw new ArgumentException("You MUST set ClientSecret and DataStore on the initializer");
      }
    }

    #region IAuthorizationCodeFlow overrides

    public Uri MakeAuthorizationCodeRequest(string state)
    {
      var v = new AuthorizationCodeRequest()
      {
        ClientId = _.ClientSecrets.ClientId,
        Scope = string.Join(" ", _.Scopes),
        State = state,
        RedirectUri = _.RedirectUrl,
      };
      return v.Build(new Uri(_.AuthorizationServerUrl));
    }

    public async Task<TokenResponse> ExchangeCodeForTokenAsync(string userId, string code, CancellationToken taskCancellationToken)
    {
      var authorizationCodeTokenReq = new AuthorizationCodeTokenRequest
      {
        Scope = string.Join(" ", _.Scopes),
        RedirectUri = _.RedirectUrl,
        Code = code,
      };

      var token = await FetchTokenAsync(userId, authorizationCodeTokenReq, new Uri(_.TokenServerUrl),  taskCancellationToken)
        .ConfigureAwait(false);
      await StoreTokenAsync(userId, token, taskCancellationToken).ConfigureAwait(false);
      return token;
    }

    public async Task<TokenResponse> RefreshTokenAsync(string userId, string refreshToken,
      CancellationToken taskCancellationToken)
    {
      var refershTokenReq = new RefreshTokenRequest
      {
        RefreshToken = refreshToken,
      };
      var token = await FetchTokenAsync(userId, refershTokenReq, new Uri( _.TokenServerUrl), taskCancellationToken).ConfigureAwait(false);

      // The new token may not contain a refresh token, so set it with the given refresh token.
      if (token.RefreshToken == null)
      {
        token.RefreshToken = refreshToken;
      }

      await StoreTokenAsync(userId, token, taskCancellationToken).ConfigureAwait(false);
      return token;
    }

    #endregion

    /// <summary>Stores the token in the DataStore/>.</summary>
    /// <param name="userId">User identifier</param>
    /// <param name="token">Token to store</param>
    /// <param name="taskCancellationToken">Cancellation token to cancel operation</param>
    private async Task StoreTokenAsync(string userId, TokenResponse token, CancellationToken taskCancellationToken)
    {
      taskCancellationToken.ThrowIfCancellationRequested();
      if (_.DataStore != null)
      {
        await _.DataStore.StoreAsync(userId, token).ConfigureAwait(false);
      }
    }

    /// <summary>Retrieve a new token from the server using the specified request.</summary>
    /// <param name="userId">User identifier</param>
    /// <param name="request">Token request</param>
    /// <param name="taskCancellationToken">Cancellation token to cancel operation</param>
    /// <returns>Token response with the new access token</returns>
    internal async Task<TokenResponse> FetchTokenAsync(string userId, TokenRequestBase request, Uri baseurl,
      CancellationToken taskCancellationToken)
    {
      // Add client id and client secret to requests.
      request.ClientId = _.ClientSecrets.ClientId;
      request.ClientSecret = _.ClientSecrets.ClientSecret;

      var rb = new RequestBuilder(request) { BaseUri = baseurl };
      var content = rb.CreateFormUrlEncodedContent();
      using (var client = new HttpClient() { BaseAddress = baseurl })
      {
        var resp = await client.PostAsync(rb.RelativePath(), content, taskCancellationToken);
        if (!resp.IsSuccessStatusCode)
        {
          var e = await resp.Content.ReadAsAsync<TokenErrorResponse>(taskCancellationToken);
          await DeleteTokenAsync(userId, taskCancellationToken).ConfigureAwait(false);
          throw new ResponseException<TokenErrorResponse>(e);
        }
        var token = await resp.Content.ReadAsAsync<TokenResponse>(taskCancellationToken);
        token.Issued = DateTime.Now;
        return token;
      }
    }

    public void Dispose()
    {
    }


    public async Task<TokenResponse> LoadTokenAsync(string userId, CancellationToken taskCancellationToken)
    {
      taskCancellationToken.ThrowIfCancellationRequested();
      if (_.DataStore == null)
      {
        return null;
      }
      return await _.DataStore.GetAsync<TokenResponse>(userId).ConfigureAwait(false);
    }

    public async Task DeleteTokenAsync(string userId, CancellationToken taskCancellationToken)
    {
      taskCancellationToken.ThrowIfCancellationRequested();
      if (_.DataStore != null)
      {
        await _.DataStore.DeleteAsync<TokenResponse>(userId).ConfigureAwait(false);
      }
    }

  }
}