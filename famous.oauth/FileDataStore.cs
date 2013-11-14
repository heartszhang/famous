using System;
using System.IO;
using System.Threading.Tasks;
using famous.oauth.core;
using Newtonsoft.Json;

namespace famous.oauth
{
  /// <summary>
  /// File data store that implements <seealso cref="IDataStorage"/>. This store creates a different file for each 
  /// combination of type and key. This file data store stores a JSON format of the specified object.
  /// </summary>
  public class FileDataStore : IDataStorage
  {
    readonly string folder_path;
    /// <summary>Gets the full folder path.</summary>
    public string FolderPath { get { return folder_path; } }

    /// <summary>
    /// Constructs a new file data store with the specified folder. This folder is created (if it doesn't exist 
    /// yet) under <seealso cref="Environment.SpecialFolder.ApplicationData"/>.
    /// </summary>
    /// <param name="folder">Folder name</param>
    public FileDataStore(string folder)
    {
      folder_path = Path.Combine(Environment.GetFolderPath(Environment.SpecialFolder.ApplicationData), folder);
      if (!Directory.Exists(folder_path))
      {
        Directory.CreateDirectory(folder_path);
      }
    }

    /// <summary>
    /// Stores the given value for the given key. It creates a new file (named <see cref="GenerateStoredKey"/>) in 
    /// <see cref="FolderPath"/>.
    /// </summary>
    /// <typeparam name="T">The type to store in the data store</typeparam>
    /// <param name="key">The key</param>
    /// <param name="value">The value to store in the data store</param>
    public Task StoreAsync<T>(string key, T value)
    {
      if (string.IsNullOrEmpty(key))
      {
        throw new ArgumentException("Key MUST have a value");
      }

      var serialized = JsonConvert.SerializeObject(value);
      var filePath = Path.Combine(folder_path, GenerateStoredKey(key, typeof(T)));
      using (var writer = File.CreateText(filePath))
      {
        return writer.WriteAsync(serialized);
      }
    }

    /// <summary>
    /// Deletes the given key. It deletes the <see cref="GenerateStoredKey"/> named file in <see cref="FolderPath"/>.
    /// </summary>
    /// <param name="key">The key to delete from the data store</param>
    public Task DeleteAsync<T>(string key)
    {
      if (string.IsNullOrEmpty(key))
      {
        throw new ArgumentException("Key MUST have a value");
      }

      var filePath = Path.Combine(folder_path, GenerateStoredKey(key, typeof(T)));
      if (File.Exists(filePath))
      {
        File.Delete(filePath);
      }
      return Task.Delay(0);
    }

    /// <summary>
    /// Returns the stored value for the given key or <c>null</c> if the matching file (<see cref="GenerateStoredKey"/>
    /// in <see cref="FolderPath"/> doesn't exist.
    /// </summary>
    /// <typeparam name="T">The type to retrieve</typeparam>
    /// <param name="key">The key to retrieve from the data store</param>
    /// <returns>The stored object</returns>
    public Task<T> GetAsync<T>(string key)
    {
      if (string.IsNullOrEmpty(key))
      {
        throw new ArgumentException("Key MUST have a value");
      }

      var tcs = new TaskCompletionSource<T>();
      var filePath = Path.Combine(folder_path, GenerateStoredKey(key, typeof(T)));
      try
      {
        if (File.Exists(filePath))
        {
          using (var reader = File.OpenText(filePath))
          {
            var re = new JsonTextReader(reader);
            var s = new JsonSerializer();
            var v = s.Deserialize<T>(re);
            tcs.SetResult(v);
          }
        }
        else
        {
          tcs.SetResult(default(T));
        }
      }
      catch (Exception e)
      {
        tcs.SetException(e);
      }
      return tcs.Task;
    }

    /// <summary>
    /// Clears all values in the data store. This method deletes all files in <see cref="FolderPath"/>.
    /// </summary>
    public Task ClearAsync()
    {
      if (!Directory.Exists(folder_path)) return Task.Delay(0);
      Directory.Delete(folder_path, true);
      Directory.CreateDirectory(folder_path);

      return Task.Delay(0);
    }

    /// <summary>Creates a unique stored key based on the key and the class type.</summary>
    /// <param name="key">The object key</param>
    /// <param name="t">The type to store or retrieve</param>
    public static string GenerateStoredKey(string key, Type t)
    {
      return string.Format("{0}-{1}", t.FullName, key);
    }
  }
}