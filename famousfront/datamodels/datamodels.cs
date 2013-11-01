using System.Runtime.Serialization;

namespace famousfront.datamodels
{
internal class FeedMediaTypes
{
  public const uint Feed_media_type_none = 0;
  public const uint Feed_media_type_unknown = 1;
  public const uint Feed_media_type_url = 1 << 1;
  public const uint Feed_media_type_video = 1 << 2;
  public const uint Feed_media_type_audio = 1 << 3;
  public const uint Feed_media_type_image = 1 << 4;
  public const uint Feed_media_type_media = Feed_media_type_audio | Feed_media_type_video;
}
internal class FeedFlags
{
  public const uint Feed_flag_none = 0;
  public const uint Feed_flag_readed = 1;
  public const uint Feed_flag_star = 1 << 1;
  public const uint Feed_flag_save = 1 << 2;
}
internal class FeedStatuses
{
  public const ulong Feed_status_fulltext_ready = 1 << 0;
  public const ulong Feed_status_fulltext_inline = 1 << 1;
  public const ulong Feed_status_thumbnail_ready = 1 << 2;
  public const ulong Feed_status_has_image = 1 << 3;
  public const ulong Feed_status_has_audio = 1 << 4;
  public const ulong Feed_status_has_video = 1 << 5;
  public const ulong Feed_status_has_url = 1 << 6;
  public const ulong Feed_status_invisible = 1 << 7;
  public const ulong Feed_status_text_only = 1 << 8;
  public const ulong Feed_status_image_only = 1 << 9;
  public const ulong Feed_status_image_gallery_only = 1 << 10;
  public const ulong Feed_status_video_only = 1 << 11;
  public const ulong Feed_status_audio_only = 1 << 12;
  public const ulong Feed_status_mp4 = 1 << 13;
  public const ulong Feed_status_flv = 1 << 14;
  public const ulong Feed_content_unresolved = 1 << 15;
  public const ulong Feed_content_ready = 1 << 16;
  public const ulong Feed_content_failed = 1 << 17;
  public const ulong Feed_content_unavail = 1 << 18;
  public const ulong Feed_content_summary = 1 << 19;
  public const ulong Feed_content_summary_replaced_by_fulltext = 1 << 20;
}
internal class FeedTypes
{
  public const uint Feed_type_unknown = 0;
  public const uint Feed_type_rss = 1 << 0;
  public const uint Feed_type_atom = 1 << 1;
  public const uint Feed_type_sinaweibo = 1 << 2;
  public const uint Feed_type_qqweibo = 1 << 3;
  public const uint Feed_type_blog = 1 << 4;
  public const uint Feed_type_tweet = 1 << 5;
  public const uint Feed_type_feed = Feed_type_rss | Feed_type_atom;
}
[DataContract]
internal class FeedLink
{
  [DataMember(EmitDefaultValue = false)]
  public string uri
  {
    get;set;
  }
  [DataMember(EmitDefaultValue = false)]
  public string alias
  {
    get;set;
  }
  [DataMember(EmitDefaultValue = false)]
  public string local
  {
    get;set;
  }
  [DataMember(EmitDefaultValue = false)]
  public string cleaned_local
  {
    get;set;
  }
  [DataMember(EmitDefaultValue = false)]
  public int words
  {
    get;set;
  }
  [DataMember(EmitDefaultValue = false)]
  public int sentences
  {
    get;set;
  }
  [DataMember(EmitDefaultValue = false)]
  public int links
  {
    get;set;
  }
  [DataMember(EmitDefaultValue = false)]
  public int density
  {
    get;set;
  }
  [DataMember(EmitDefaultValue = false)]
  public long length
  {
    get;set;
  }
  [DataMember(EmitDefaultValue = false)]
  public bool readable
  {
    get;set;
  }
  [DataMember(EmitDefaultValue = false)]
  public FeedMedia[] images
  {
    get;set;
  }
  [DataMember(EmitDefaultValue = false)]
  public FeedMedia[] videos
  {
    get;set;
  }
  [DataMember(EmitDefaultValue = false)]
  public FeedMedia[] audios
  {
    get;set;
  }
}
[DataContract]
internal class FeedMedia
{
  [DataMember(EmitDefaultValue = false)]
  public string title
  {
    get;set;
  }
  [DataMember(EmitDefaultValue = false)]
  public string description
  {
    get;set;
  }
  [DataMember(EmitDefaultValue = false)]
  public string uri
  {
    get;set;
  }
  [DataMember(EmitDefaultValue = false)]
  public string local
  {
    get;set;
  }
  [DataMember(EmitDefaultValue = false)]
  public int width
  {
    get;set;
  }
  [DataMember(EmitDefaultValue = false)]
  public int height
  {
    get;set;
  }
  [DataMember(EmitDefaultValue = false)]
  public int duration
  {
    get;set;
  }
  [DataMember(EmitDefaultValue = false)]
  public string mime
  {
    get;set;
  }
  [DataMember(EmitDefaultValue = false)]
  public string thumbanil
  {
    get;set;
  }
}
[DataContract]
internal class FeedAuthor
{
  [DataMember(EmitDefaultValue = false)]
  public string name
  {
    get;set;
  }
  [DataMember(EmitDefaultValue = false)]
  public string email
  {
    get;set;
  }
  [DataMember(EmitDefaultValue = false)]
  public ulong id
  {
    get;set;
  }
}
[DataContract]
internal class FeedTitle
{
  [DataMember(EmitDefaultValue = false)]
  public string main
  {
    get;set;
  }
  [DataMember(EmitDefaultValue = false)]
  public string[] others
  {
    get;set;
  }
}
[DataContract]
internal class FeedContent
{
  [DataMember(EmitDefaultValue = false)]
  public string uri
  {
    get;set;
  }
  [DataMember(EmitDefaultValue = false)]
  public string local
  {
    get;set;
  }
  [DataMember(EmitDefaultValue = false)]
  public uint words
  {
    get;set;
  }
  [DataMember(EmitDefaultValue = false)]
  public uint density
  {
    get;set;
  }
  [DataMember(EmitDefaultValue = false)]
  public uint links
  {
    get;set;
  }
  [DataMember(EmitDefaultValue = false)]
  public ulong status
  {
    get;set;
  }
  [DataMember(EmitDefaultValue = false)]
  public FeedMedia[] images
  {
    get;set;
  }
  [DataMember(EmitDefaultValue = false)]
  public FeedMedia[] medias
  {
    get;set;
  }
}
[DataContract]
internal class FeedEntry
{
  [DataMember(Name = "_id")]
  public string id
  {
    get;set;
  }
  [DataMember(EmitDefaultValue = false)]
  public uint flags
  {
    get;set;
  }
  [DataMember(EmitDefaultValue = false)]
  public string source
  {
    get;set;
  }
  [DataMember(EmitDefaultValue = false)]
  public uint type
  {
    get;set;
  }
  [DataMember(EmitDefaultValue = false)]
  public string uri
  {
    get;set;
  }
  [DataMember(EmitDefaultValue = false)]
  public FeedTitle title
  {
    get;set;
  }
  [DataMember(EmitDefaultValue = false)]
  public FeedAuthor author
  {
    get;set;
  }
  [DataMember(EmitDefaultValue = false)]
  public long pubdate
  {
    get;set;
  }
  [DataMember(EmitDefaultValue = false)]
  public string summary
  {
    get;set;
  }
  [DataMember(EmitDefaultValue = false)]
  public FeedContent content
  {
    get;set;
  }
  [DataMember(EmitDefaultValue = false)]
  public string[]tags
  {
    get;set;
  }
  [DataMember(EmitDefaultValue = false)]
  public FeedMedia[] images
  {
    get;set;
  }
  [DataMember(EmitDefaultValue = false)]
  public FeedMedia[]videos
  {
    get;set;
  }
  [DataMember(EmitDefaultValue = false)]
  public FeedMedia[]audios
  {
    get;set;
  }
  [DataMember(EmitDefaultValue = false)]
  public FeedLink[] links
  {
    get;set;
  }
  [DataMember(EmitDefaultValue = false)]
  public uint words
  {
    get;set;
  }
  [DataMember(EmitDefaultValue = false)]
  public uint density
  {
    get;set;
  }
  [DataMember(EmitDefaultValue = false)]
  public ulong status
  {
    get;set;
  }
  [DataMember(EmitDefaultValue = false)]
  public string[] categories
  {
    get;set;
  }
}
[DataContract]
internal class FeedSource
{
  [DataMember(EmitDefaultValue = false)]
  public string name
  {
    get;set;
  }
  [DataMember(EmitDefaultValue = false)]
  public string uri
  {
    get;set;
  }
  [DataMember(EmitDefaultValue = false)]
  public string local
  {
    get;set;
  }
  [DataMember(EmitDefaultValue = false)]
  public uint period
  {
    get;set;
  }
  [DataMember(EmitDefaultValue = false)]
  public ulong deadline
  {
    get;set;
  }
  [DataMember(EmitDefaultValue = false)]
  public uint type
  {
    get;set;
  }
  [DataMember(EmitDefaultValue = false)]
  public bool disabled
  {
    get;set;
  }
  [DataMember(EmitDefaultValue = false)]
  public bool enable_proxy
  {
    get;set;
  }
  [DataMember(EmitDefaultValue = false)]
  public ulong update
  {
    get;set;
  }
  [DataMember(EmitDefaultValue = false)]
  public string website
  {
    get;set;
  }
  [DataMember(EmitDefaultValue = false)]
  public string[] tags
  {
    get;set;
  }
  [DataMember(EmitDefaultValue = false)]
  public string[] categories
  {
    get;set;
  }
  [DataMember(EmitDefaultValue = false)]
  public int unreaded
  {
    get;set;
  }
  [DataMember(EmitDefaultValue = false)]
  public FeedMedia media
  {
    get;set;
  }
  [DataMember(EmitDefaultValue = false)]
  public string description
  {
    get;set;
  }
}
[DataContract]
internal class FeedTag
{
  [DataMember(Name = "value")]
  public string Name
  {
    get;set;
  }
}
internal class FeedSourceTypes
{
  public const string Feed_type_unknown = "";
  public const string Feed_type_rss = "rss";
  public const string Feed_type_atom = "atom";
  public const string Feed_type_blog = "blog";
  public const string Feed_type_tweet = "tweet";
  public const string Feed_type_sinaweibo = "weibo";
  public const string Feed_type_qqweibo = "qqweibo";
}
}