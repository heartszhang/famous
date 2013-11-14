using System;
using System.Collections.Generic;
using System.Linq;
using System.Text;
using System.Threading;
using System.Threading.Tasks;
using famous.oauth.requests;
using famous.oauth.responses;

namespace famous.oauth.core
{
    /// <summary>OAuth 2.0 authorization code flow that manages and persists end-user credentials.</summary>
  public interface IAuthorizationCodeFlow : IDisposable
  {
    /// <summary>Asynchronously loads the user's token using the flow's <seealso cref="IDataStorage"/>.</summary>
    /// <param name="userId">User identifier</param>
    /// <param name="taskCancellationToken">Cancellation token to cancel operation</param>
    /// <returns>Token response</returns>
    Task<TokenResponse> LoadTokenAsync(string userId, CancellationToken taskCancellationToken);

    /// <summary>Asynchronously deletes the user's token using the flow's <seealso cref="IDataStorage"/>.</summary>
    /// <param name="userId">User identifier</param>
    /// <param name="taskCancellationToken">Cancellation token to cancel operation</param>
    Task DeleteTokenAsync(string userId, CancellationToken taskCancellationToken);

    /// <summary>Creates an authorization code request with the specified redirect URI.</summary>
    Uri MakeAuthorizationCodeRequest(string state);

    /// <summary>Asynchronously exchanges code with a token.</summary>
    /// <param name="userId">User identifier</param>
    /// <param name="code">Authorization code received from the authorization server</param>
    /// <param name="taskCancellationToken">Cancellation token to cancel operation</param>
    /// <returns>Token response which contains the access token</returns>
    Task<TokenResponse> ExchangeCodeForTokenAsync(string userId, string code, CancellationToken taskCancellationToken);

    /// <summary>Asynchronously refreshes an access token using a refresh token.</summary>
    /// <param name="userId">User identifier</param>
    /// <param name="refreshToken">Refresh token which is used to get a new access token</param>
    /// <param name="taskCancellationToken">Cancellation token to cancel operation</param>
    /// <returns>Token response which contains the access token and the input refresh token</returns>
    Task<TokenResponse> RefreshTokenAsync(string userId, string refreshToken,
        CancellationToken taskCancellationToken);
  }
}
